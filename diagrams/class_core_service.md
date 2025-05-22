````mermaid
classDiagram
    class Comment {
    -int commentId
    -int eventId
    -int senderId
    -int taskId
    -string content
    -DateTime createdAt
    -bool isDeleted
    -bool isRead
    +createComment() : Comment
    +listByEvent() : List~Comment~
    }

    class User {
    -int userId
    -string username
    -string email
    -string role
    -bool isActive
    -bool isDeleted
    -DateTime createdAt
    +getUserById() : User
    +getUserByUsername() : User
    +listUsers() : List~User~
    }

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
    +createEvent() : Event
    +updateEvent() : void
    +deleteEvent() : void
    +getEventDetails() : Event
    }

    class Task {
    -int taskId
    -int eventId
    -int parentId
    -string title
    -string description
    -int storyPoints
    -string priority
    -string status
    -DateTime createdAt
    +createTask() : Task
    +updateTaskStatus() : void
    +getTaskById() : Task
    +listByEvent() : List~Task~
    }

    class TaskAssignment {
    -int taskAssignmentId
    -int taskId
    -int userId
    -DateTime assignedAt
    -DateTime completedAt
    +assignTask() : void
    +markTaskCompleted() : void
    }

    class Expense {
    -int expenseId
    -int eventId
    -int taskId
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
    -float amount
    -boolean isPaid
    -DateTime paidAt
    +updatePaymentStatus() : void
    }

    class RabbitMQPublisher {
    -string queueName
    -Connection connection
    -Channel channel
    +publish(data : byte[]) : void
    +stop() : void
    }

    class EventParticipant {
    -int eventParticipantId
    -int eventId
    -int userId
    -string role
    -DateTime joinedAt
    -bool isConfirmed
    +confirmParticipation() : void
    +listByEvent() : List~EventParticipant~
    +listByUser() : List~EventParticipant~
    }

    Event *--> EventParticipant : has
    EventParticipant o-- User : participant
    Event *--> Task : contains
    Event *--> Expense : includes
    Event o-- User : organized_by
    Task *--> TaskAssignment : assigned_to
    Expense *--> ExpenseShare : splits
    User o-- TaskAssignment : assigned
    User o-- ExpenseShare : owes
    Event --> RabbitMQPublisher : sends
    Task --> RabbitMQPublisher : sends
    Comment o-- Event : belongs_to
    Comment o-- Task : optional_for
    EventParticipant --> RabbitMQPublisher : sends
````