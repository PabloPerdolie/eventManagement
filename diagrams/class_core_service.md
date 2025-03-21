````mermaid
classDiagram
    class Event {
        -int eventId
        -int organizerId
        -string title
        -string description
        -DateTime startDate
        -DateTime endDate
        -string location
        -string status
        -DateTime createdAt
        -DateTime updatedAt
        +createEvent() : Event
        +updateEvent() : void
        +deleteEvent() : void
        +getEventDetails() : Event
    }

    class User {
        -int userId
        -string username
        -string email
        -string role
        +getUserById() : User
        +addParticipant(eventId : int) : void
    }

    class Task {
        -int taskId
        -int eventId
        -string title
        -string description
        -DateTime dueDate
        -string priority
        -string status
        -DateTime createdAt
        -DateTime updatedAt
        +createTask() : Task
        +updateTaskStatus(status : string) : void
        +getTaskById() : Task
    }

    class TaskAssignment {
        -int assignmentId
        -int taskId
        -int userId
        -DateTime assignedAt
        -DateTime completedAt
        +assignTask(userId : int) : void
        +markTaskCompleted() : void
    }

    class Expense {
        -int expenseId
        -int eventId
        -string description
        -float amount
        -string currency
        -int createdBy
        -string splitMethod
        -DateTime createdAt
        +addExpense() : Expense
        +calculateShares() : List~ExpenseShare~
    }

    class ExpenseShare {
        -int shareId
        -int expenseId
        -int userId
        -float shareAmount
        -boolean isPaid
        -DateTime paidAt
        +updatePaymentStatus(isPaid : boolean) : void
    }

    class NotificationPublisher {
        -string queueName
        -RabbitMQConnection connection
        +publishNotification(event : string, data : Object) : void
    }

    class RabbitMQConnection {
        -string host
        -int port
        -string username
        -string password
        +connect() : void
        +disconnect() : void
        +getChannel() : Channel
    }

    class EventParticipant {
        -int eventParticipantId
        -int eventId
        -int userId
        -string status
        +confirmParticipation() : void
    }

    Event *--> EventParticipant : has
EventParticipant o-- User : participant
    Task --> NotificationPublisher: sends
    Event *--> Task : contains
    Event *--> Expense : includes
    Event o-- User : organized_by
    Task *--> TaskAssignment : assigned_to
    Expense *--> ExpenseShare : splits
    User o-- TaskAssignment : assigned
    User o-- ExpenseShare : owes
    NotificationPublisher --> RabbitMQConnection : uses
````