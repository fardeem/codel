package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/semanser/ai-coder/executor"
	gmodel "github.com/semanser/ai-coder/graph/model"
	"github.com/semanser/ai-coder/graph/subscriptions"
	"github.com/semanser/ai-coder/models"
	"github.com/semanser/ai-coder/services"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CreateFlow is the resolver for the createFlow field.
func (r *mutationResolver) CreateFlow(ctx context.Context, query string) (*gmodel.Flow, error) {
	summary, err := services.GetMessageSummary(query, 10)

	if err != nil {
		return nil, fmt.Errorf("failed to get message summary: %w", err)
	}

	dockerImage, err := services.GetDockerImageName(query)

	if err != nil {
		return nil, fmt.Errorf("failed to get docker image name: %w", err)
	}

	flow := models.Flow{
		Name:        summary,
		DockerImage: dockerImage,
	}

	tx := r.Db.Create(&flow)

	if tx.Error != nil {
		return nil, tx.Error
	}

	_, err = executor.SpawnContainer(executor.GenerateContainerName(flow.ID), dockerImage)

	if err != nil {
		return nil, fmt.Errorf("failed to spawn container: %w", err)
	}

	task, err := r.CreateTask(ctx, flow.ID, query)

	if err != nil {
		return nil, fmt.Errorf("failed to create the inital task: %w", err)
	}

	return &gmodel.Flow{
		ID:    flow.ID,
		Name:  flow.Name,
		Tasks: []*gmodel.Task{task},
	}, nil
}

// CreateTask is the resolver for the createTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, id uint, query string) (*gmodel.Task, error) {
	type InputTaskArgs struct {
		Query string `json:"query"`
	}

	args := InputTaskArgs{Query: query}
	arg, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	flowResult := r.Db.First(&models.Flow{}, id)

	if errors.Is(flowResult.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("flow with id %d not found", id)
	}

	task := models.Task{
		Type:    models.Input,
		Message: query,
		Status:  models.Finished,
		Args:    datatypes.JSON(arg),
		FlowID:  id,
	}

	tx := r.Db.Create(&task)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to create task: %w", tx.Error)
	}

	executor.AddCommand(task)

	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	return &gmodel.Task{
		ID:        task.ID,
		Message:   task.Message,
		Type:      gmodel.TaskType(task.Type),
		Status:    gmodel.TaskStatus(task.Status),
		Args:      task.Args.String(),
		CreatedAt: task.CreatedAt,
	}, nil
}

// StopTask is the resolver for the stopTask field.
func (r *mutationResolver) StopTask(ctx context.Context, id uint) (*gmodel.Task, error) {
	panic(fmt.Errorf("not implemented: StopTask - stopTask"))
}

// Exec is the resolver for the _exec field.
func (r *mutationResolver) Exec(ctx context.Context, containerID string, command string) (string, error) {
	b := bytes.Buffer{}
	executor.ExecCommand(containerID, command, &b)

	return b.String(), nil
}

// Flows is the resolver for the flows field.
func (r *queryResolver) Flows(ctx context.Context) ([]*gmodel.Flow, error) {
	flows := []models.Flow{}

	tx := r.Db.Model(&models.Flow{}).Order("created_at DESC").Preload("Tasks").Find(&flows)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to fetch flows: %w", tx.Error)
	}

	var gFlows []*gmodel.Flow

	for _, flow := range flows {
		var gTasks []*gmodel.Task

		for _, task := range flow.Tasks {
			gTasks = append(gTasks, &gmodel.Task{
				ID:        task.ID,
				Message:   task.Message,
				Type:      gmodel.TaskType(task.Type),
				Status:    gmodel.TaskStatus(task.Status),
				Args:      task.Args.String(),
				Results:   task.Results,
				CreatedAt: task.CreatedAt,
			})

		}

		gFlows = append(gFlows, &gmodel.Flow{
			ID:    flow.ID,
			Name:  flow.Name,
			Tasks: gTasks,
		})
	}

	return gFlows, nil
}

// Flow is the resolver for the flow field.
func (r *queryResolver) Flow(ctx context.Context, id uint) (*gmodel.Flow, error) {
	flow := models.Flow{}

	tx := r.Db.First(&models.Flow{}, id).Preload("Tasks").Find(&flow)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to fetch flows: %w", tx.Error)
	}

	var gFlow *gmodel.Flow
	var gTasks []*gmodel.Task

	for _, task := range flow.Tasks {
		gTasks = append(gTasks, &gmodel.Task{
			ID:        task.ID,
			Message:   task.Message,
			Type:      gmodel.TaskType(task.Type),
			Status:    gmodel.TaskStatus(task.Status),
			Args:      task.Args.String(),
			Results:   task.Results,
			CreatedAt: task.CreatedAt,
		})
	}

	gFlow = &gmodel.Flow{
		ID:    flow.ID,
		Name:  flow.Name,
		Tasks: gTasks,
	}

	return gFlow, nil
}

// TaskAdded is the resolver for the taskAdded field.
func (r *subscriptionResolver) TaskAdded(ctx context.Context, flowID uint) (<-chan *gmodel.Task, error) {
	return subscriptions.TaskAdded(ctx, flowID)
}

// TaskUpdated is the resolver for the taskUpdated field.
func (r *subscriptionResolver) TaskUpdated(ctx context.Context) (<-chan *gmodel.Task, error) {
	panic(fmt.Errorf("not implemented: TaskUpdated - taskUpdated"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
