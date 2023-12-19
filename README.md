# Go Blocking Queue

## Introduction

This documentation covers a Go code implementation of a multi-threaded task processing system using blocking queues.

## Code Structure

The code consists of the following main components:

### `Task` Struct

The `Task` struct represents a task with various fields, including task type (`Typ`), start and end times (`StartTime` and `EndTime`), an ID, data, status, result, error, owner, and timestamps for creation and last update.

### `BlockingQueue` Struct

The `BlockingQueue` struct is a blocking queue implementation that allows task insertion and retrieval. It utilizes a mutex and conditional variable for synchronization.

### `NewBlockingQueue` Function

`NewBlockingQueue` creates and initializes a new `BlockingQueue` instance.

### `Put` Method

The `Put` method adds a task to the queue. It locks the queue, appends the task to the data slice, and signals waiting consumers. It also sets the `StartTime` and `CreatedAt` fields of the task.

### `Get` Method

The `Get` method retrieves and removes a task from the queue. If the queue is empty, it blocks until a task is available. It also sets the `EndTime` and `UpdatedAt` fields of the task.

### `Size` Method

The `Size` method returns the current number of tasks in the queue.

### `Close` Method

The `Close` method marks the queue as closed, allowing consumers to gracefully exit.

## Exemple usage


   ```go
    package main

    import (
        "fmt"
        "time"

        "blocking_queue" 
    )

    func main() {
        // Create a new BlockingQueue instance
        queue := blocking_queue.NewBlockingQueue()

        // Create a task and put it into the queue
        task := &blocking_queue.Task{
            Typ:       'A',             // Task type
            StartTime: time.Now(),      // Record start time
            EndTime:   time.Time{},     // Set to zero value initially
        }
        if err := queue.Put(task); err != nil {
            fmt.Printf("Error putting task into the queue: %v\n", err)
            return
        }

        // Retrieve a task from the queue
        retrievedTask, err := queue.Get()
        if err != nil {
            fmt.Printf("Error getting task from the queue: %v\n", err)
            return
        }

        // Check the size of the queue
        size := queue.Size()

        fmt.Printf("Task retrieved from the queue: %+v\n", retrievedTask)
        fmt.Printf("Current size of the queue: %d\n", size)
    }
