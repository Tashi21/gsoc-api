package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	Id        uuid.UUID  `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Year      int        `json:"year,omitempty"`
	TechStack *string    `json:"tech_stack,omitempty"`
	Topics    *string    `json:"topics,omitempty"`
	ShortDesc *string    `json:"short_desc,omitempty"`
	Link      *string    `json:"link,omitempty"`
	ImgUrl    *string    `json:"img_url,omitempty"`
	Website   *string    `json:"website,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type JsonResponse struct {
	Type    string         `json:"type,omitempty"`
	Message string         `json:"message,omitempty"`
	Count   int            `json:"count,omitempty"`
	Data    []Organization `json:"data,omitempty"`
}
