package controller

import "ratblog/contract"

type controller struct {
	repo contract.Repository
}

func New(repo contract.Repository) *controller {
	return &controller{repo: repo}
}
