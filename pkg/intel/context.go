package intel

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ContextManager handles intelligent context management with token awareness
type ContextManager struct {
	maxTokens     int
	currentTokens int
	relevanceDecay float64
	items         []ContextItem
}

// ContextItem represents a piece of context with metadata
type ContextItem struct {
	ID          string
	Type        ContextType
	Content     string
	Timestamp   time.Time
	Relevance   float64
	TokenCount  int
	IsEssential bool
}

// ContextType defines different types of context
type ContextType int

const (
	ContextTypeSystem ContextType = iota
	ContextTypeDomain
	ContextTypeState
	ContextTypeHistory
	ContextTypePrompt
	ContextTypeUser
)

// String returns the string representation of ContextType
func (ct ContextType) String() string {
	switch ct {
	case ContextTypeSystem:
		return "system"
	case ContextTypeDomain:
		return "domain"
	case ContextTypeState:
		return "state"
	case ContextTypeHistory:
		return "history"
	case ContextTypePrompt:
		return "prompt"
	case ContextTypeUser:
		return "user"
	default:
		return "unknown"
	}
}

// NewContextManager creates a new context manager
func NewContextManager(maxTokens int) *ContextManager {
	return &ContextManager{
		maxTokens:     maxTokens,
		currentTokens: 0,
		relevanceDecay: 0.9, // Decay factor for aging context
		items:         make([]ContextItem, 0),
	}
}

// estimateTokens provides a rough estimate of token count for text
func (cm *ContextManager) estimateTokens(text string) int {
	// Rough estimation: ~4 characters per token for English text
	return len(text) / 4
}

// AddContext adds a new context item
func (cm *ContextManager) AddContext(id string, contextType ContextType, content string, isEssential bool) {
	tokens := cm.estimateTokens(content)
	
	item := ContextItem{
		ID:          id,
		Type:        contextType,
		Content:     content,
		Timestamp:   time.Now(),
		Relevance:   cm.calculateInitialRelevance(contextType),
		TokenCount:  tokens,
		IsEssential: isEssential,
	}
	
	// Remove existing item with same ID
	cm.removeItem(id)
	
	// Add new item
	cm.items = append(cm.items, item)
	cm.currentTokens += tokens
	
	// Optimize if we exceed limits
	if cm.currentTokens > cm.maxTokens {
		cm.optimize()
	}
}

// calculateInitialRelevance calculates initial relevance based on context type
func (cm *ContextManager) calculateInitialRelevance(contextType ContextType) float64 {
	switch contextType {
	case ContextTypeSystem:
		return 1.0 // System prompts are always relevant
	case ContextTypeDomain:
		return 0.9 // Domain knowledge is very relevant
	case ContextTypeState:
		return 0.8 // Current state is important
	case ContextTypePrompt:
		return 0.7 // Custom prompts are relevant
	case ContextTypeHistory:
		return 0.6 // History is somewhat relevant
	case ContextTypeUser:
		return 0.9 // User queries are very relevant
	default:
		return 0.5 // Default relevance
	}
}

// removeItem removes an item with the given ID
func (cm *ContextManager) removeItem(id string) {
	for i, item := range cm.items {
		if item.ID == id {
			cm.currentTokens -= item.TokenCount
			cm.items = append(cm.items[:i], cm.items[i+1:]...)
			break
		}
	}
}

// optimize reduces context size by removing less relevant items
func (cm *ContextManager) optimize() {
	// Update relevance scores based on age
	cm.updateRelevanceScores()
	
	// Sort by relevance (essential items first, then by relevance score)
	sort.Slice(cm.items, func(i, j int) bool {
		if cm.items[i].IsEssential != cm.items[j].IsEssential {
			return cm.items[i].IsEssential
		}
		return cm.items[i].Relevance > cm.items[j].Relevance
	})
	
	// Remove items until we're under the limit
	targetTokens := int(float64(cm.maxTokens) * 0.8) // Leave 20% buffer
	
	for cm.currentTokens > targetTokens && len(cm.items) > 0 {
		// Find the least relevant non-essential item
		var removeIndex = -1
		for i := len(cm.items) - 1; i >= 0; i-- {
			if !cm.items[i].IsEssential {
				removeIndex = i
				break
			}
		}
		
		if removeIndex == -1 {
			// All items are essential, can't remove more
			break
		}
		
		// Remove the item
		cm.currentTokens -= cm.items[removeIndex].TokenCount
		cm.items = append(cm.items[:removeIndex], cm.items[removeIndex+1:]...)
	}
}

// updateRelevanceScores updates relevance scores based on age and usage
func (cm *ContextManager) updateRelevanceScores() {
	now := time.Now()
	
	for i := range cm.items {
		item := &cm.items[i]
		
		// Age-based decay
		age := now.Sub(item.Timestamp)
		ageFactor := 1.0
		
		// Different decay rates for different types
		switch item.Type {
		case ContextTypeHistory:
			// History decays faster
			ageFactor = 1.0 - (age.Minutes() / 60.0) * 0.5
		case ContextTypeState:
			// State decays slower
			ageFactor = 1.0 - (age.Minutes() / 120.0) * 0.3
		case ContextTypeSystem, ContextTypeDomain:
			// System and domain knowledge don't decay
			ageFactor = 1.0
		default:
			// Default decay
			ageFactor = 1.0 - (age.Minutes() / 90.0) * 0.4
		}
		
		if ageFactor < 0.1 {
			ageFactor = 0.1 // Minimum relevance
		}
		
		item.Relevance *= ageFactor
	}
}

// BuildPrompt constructs an optimized prompt from available context
func (cm *ContextManager) BuildPrompt(userQuery string, promptType PromptType) string {
	var prompt strings.Builder
	
	// Update relevance scores
	cm.updateRelevanceScores()
	
	// Sort items by type priority and relevance
	sortedItems := make([]ContextItem, len(cm.items))
	copy(sortedItems, cm.items)
	
	sort.Slice(sortedItems, func(i, j int) bool {
		// System context first
		if sortedItems[i].Type != sortedItems[j].Type {
			return cm.getTypePriority(sortedItems[i].Type) > cm.getTypePriority(sortedItems[j].Type)
		}
		return sortedItems[i].Relevance > sortedItems[j].Relevance
	})
	
	// Build prompt with context items
	for _, item := range sortedItems {
		if item.Type == ContextTypeUser {
			continue // Handle user query separately
		}
		
		prompt.WriteString(item.Content)
		prompt.WriteString("\n\n")
	}
	
	// Add user query
	prompt.WriteString("Query: ")
	prompt.WriteString(userQuery)
	prompt.WriteString("\n\nResponse:")
	
	return prompt.String()
}

// getTypePriority returns priority for context types
func (cm *ContextManager) getTypePriority(contextType ContextType) int {
	switch contextType {
	case ContextTypeSystem:
		return 100
	case ContextTypeDomain:
		return 90
	case ContextTypePrompt:
		return 80
	case ContextTypeState:
		return 70
	case ContextTypeHistory:
		return 60
	case ContextTypeUser:
		return 50
	default:
		return 40
	}
}

// GetStats returns current context statistics
func (cm *ContextManager) GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"total_items":    len(cm.items),
		"current_tokens": cm.currentTokens,
		"max_tokens":     cm.maxTokens,
		"utilization":    float64(cm.currentTokens) / float64(cm.maxTokens),
	}
	
	// Count by type
	typeCounts := make(map[string]int)
	for _, item := range cm.items {
		typeCounts[item.Type.String()]++
	}
	stats["by_type"] = typeCounts
	
	return stats
}

// Clear removes all context items
func (cm *ContextManager) Clear() {
	cm.items = make([]ContextItem, 0)
	cm.currentTokens = 0
}

// GetContextSummary returns a summary of current context
func (cm *ContextManager) GetContextSummary() string {
	if len(cm.items) == 0 {
		return "No context available"
	}
	
	var summary strings.Builder
	summary.WriteString("Current context:\n")
	
	typeCounts := make(map[ContextType]int)
	for _, item := range cm.items {
		typeCounts[item.Type]++
	}
	
	for contextType, count := range typeCounts {
		summary.WriteString("- ")
		summary.WriteString(contextType.String())
		summary.WriteString(": ")
		summary.WriteString(fmt.Sprintf("%d items", count))
		summary.WriteString("\n")
	}
	
	summary.WriteString(fmt.Sprintf("Total tokens: %d/%d (%.1f%%)", 
		cm.currentTokens, cm.maxTokens, 
		float64(cm.currentTokens)/float64(cm.maxTokens)*100))
	
	return summary.String()
}

// PruneHistory removes old history items to make space
func (cm *ContextManager) PruneHistory() {
	cutoff := time.Now().Add(-30 * time.Minute) // Remove history older than 30 minutes
	
	newItems := make([]ContextItem, 0)
	for _, item := range cm.items {
		if item.Type == ContextTypeHistory && item.Timestamp.Before(cutoff) && !item.IsEssential {
			cm.currentTokens -= item.TokenCount
			continue
		}
		newItems = append(newItems, item)
	}
	
	cm.items = newItems
}

// SetMaxTokens updates the maximum token limit
func (cm *ContextManager) SetMaxTokens(maxTokens int) {
	cm.maxTokens = maxTokens
	if cm.currentTokens > cm.maxTokens {
		cm.optimize()
	}
}

// GetMostRelevantItems returns the most relevant context items
func (cm *ContextManager) GetMostRelevantItems(count int) []ContextItem {
	if len(cm.items) == 0 {
		return []ContextItem{}
	}
	
	// Sort by relevance
	sortedItems := make([]ContextItem, len(cm.items))
	copy(sortedItems, cm.items)
	
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].Relevance > sortedItems[j].Relevance
	})
	
	if count > len(sortedItems) {
		count = len(sortedItems)
	}
	
	return sortedItems[:count]
}