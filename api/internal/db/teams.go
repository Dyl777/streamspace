package db

import "time"

// TeamMembership represents a user's membership in a team
type TeamMembership struct {
	TeamID          string    `json:"teamId"`
	TeamName        string    `json:"teamName"`
	TeamDisplayName string    `json:"teamDisplayName"`
	TeamType        string    `json:"teamType"`
	Role            string    `json:"role"`
	JoinedAt        time.Time `json:"joinedAt"`
}

// TeamPermission represents a permission for a team role
type TeamPermission struct {
	ID          int       `json:"id"`
	Role        string    `json:"role"`
	Permission  string    `json:"permission"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

// TeamRoleInfo represents information about a team role and its permissions
type TeamRoleInfo struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}
