````mermaid
erDiagram
    USER ||--o{ EXPENSE_SHARE : owes
    USER ||--o{ MESSAGE : sends
    USER ||--o{ EVENT : organizes
    USER ||--o{ EVENT_PARTICIPANT : participates
    EVENT ||--o{ TASK : contains
    EVENT ||--o{ EXPENSE : includes
    EVENT ||--o{ MESSAGE : includes
    TASK ||--o{ TASK_ASSIGNMENT : assigned_to
    EXPENSE ||--o{ EXPENSE_SHARE : splits
    USER ||--o{ TASK_ASSIGNMENT : assigned
    

    USER {
        int user_id PK
        string username
        string email
        string password_hash
        string first_name
        string last_name
        timestamp created_at
        timestamp updated_at
        boolean is_active
        string role "organizer, participant"
    }

    EVENT {
        int event_id PK
        int organizer_id FK "USER"
        string title
        string description
        timestamp start_date
        timestamp end_date
        string location
        string status "planned, active, completed"
        timestamp created_at
        timestamp updated_at
    }

    EVENT_PARTICIPANT {
        int event_participant_id PK
        int event_id FK "EVENT"
        int user_id FK "USER"
        timestamp joined_at
        boolean is_confirmed
    }

    TASK {
        int task_id PK
        int event_id FK "EVENT"
        string title
        string description
        timestamp due_date
        string priority "low, medium, high"
        string status "pending, in_progress, completed"
        timestamp created_at
        timestamp updated_at
    }

    TASK_ASSIGNMENT {
        int task_assignment_id PK
        int task_id FK "TASK"
        int user_id FK "USER"
        timestamp assigned_at
        timestamp completed_at
    }

    EXPENSE {
        int expense_id PK
        int event_id FK "EVENT"
        string description
        float amount
        string currency "USD, EUR, etc."
        timestamp expense_date
        int created_by FK "USER"
        string split_method "equal, custom"
        timestamp created_at
    }

    EXPENSE_SHARE {
        int expense_share_id PK
        int expense_id FK "EXPENSE"
        int user_id FK "USER"
        float share_amount
        boolean is_paid
        timestamp paid_at
    }

    MESSAGE {
        int message_id PK
        int event_id FK "EVENT"
        int sender_id FK "USER"
        string content
        timestamp sent_at
        boolean is_read
    }
````