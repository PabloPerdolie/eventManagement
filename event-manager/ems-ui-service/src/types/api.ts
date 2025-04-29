export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: UserResponse;
}

export interface TokenResponse {
  access_token: string;
  expires_in: number;
}

export interface UserResponse {
  id: number;
  username: string;
  email: string;
  role: string;
  created_at: string;
  is_active?: boolean;
}

export interface UserLoginRequest {
  username: string;
  password: string;
}

export interface UserRegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface SuccessResponse {
  success: boolean;
  message: string;
  data?: any;
}

export interface ErrorResponse {
  error: string;
  message: string;
}

export interface EventResponse {
  id: number;
  title: string;
  description: string;
  location: string;
  start_date: string;
  end_date: string;
  created_by: number;
  created_at: string;
  status: string;
}

export interface EventCreateRequest {
  title: string;
  description?: string;
  location?: string;
  start_date: string;
  end_date: string;
}

export interface EventsResponse {
  events: EventResponse[];
  total: number;
}

export interface TaskResponse {
  id: number;
  title: string;
  description: string;
  event_id: number;
  assigned_to: number;
  parent_id: number | null;
  priority: string;
  status: TaskStatus;
  story_points: number;
  created_at: string;
}

export interface TaskCreateRequest {
  title: string;
  description?: string;
  event_id: number;
  assigned_to?: number;
  parent_id?: number;
  priority?: string;
  story_points?: number;
}

export interface TaskUpdateRequest {
  title?: string;
  description?: string;
  assigned_to?: number;
  parent_id?: number;
  priority?: string;
  status?: TaskStatus;
  story_points?: number;
}

export interface TasksResponse {
  tasks: TaskResponse[];
  total: number;
}

export enum TaskStatus {
  Pending = 'pending',
  InProgress = 'in_progress',
  Completed = 'completed',
  Cancelled = 'cancelled'
}

export interface EventParticipant {
  id: number;
  event_id: number;
  user: UserResponse;
  role: string;
  is_confirmed: boolean;
  joined_at: string;
}

export interface EventParticipantCreateRequest {
  user_id: number;
  event_title: string;
}

export interface EventParticipantsResponse {
  participants: EventParticipant[];
  total: number;
}

export interface Comment {
  commentId: number;
  content: string;
  eventId: number;
  taskId?: number;
  senderId: number;
  createdAt: string;
  isRead: boolean;
  isDeleted: boolean;
}

export interface CommentCreateRequest {
  content: string;
  event_id: number;
  task_id?: number;
  sender_id: number;
}

export interface CommunicationServiceResponse {
  comments: Comment[];
}

export interface EventData {
  eventData: EventResponse;
  eventParticipants: EventParticipantsResponse;
  tasks: TasksResponse;
  comments: CommunicationServiceResponse;
  expenses: ExpensesResponse;
  balanceReport: BalanceReportResponse;
}

export interface Health {
  service: string;
  status: string;
  timestamp: string;
  version: string;
}

// Типы для работы с расходами
export enum SplitMethod {
  EQUAL = 'equal',
  PERCENTAGE = 'percentage',
  AMOUNT = 'amount',
  CUSTOM = 'custom'
}

export interface ExpenseShare {
  ShareID: number;
  ExpenseID: number;
  UserID: number;
  Amount: number;
  IsPaid: boolean;
  PaidAt: string | null;
}

export interface ExpenseCreateRequest {
  event_id: number;
  created_by: number;
  description: string;
  amount: number;
  currency: string;
  split_method: string;
  user_ids: number[];
}

export interface ExpenseUpdateRequest {
  description?: string;
  amount?: number;
  currency?: string;
  split_method?: string;
  user_ids?: number[];
}

export interface ExpenseShareUpdateRequest {
  is_paid: boolean;
}

export interface ExpenseResponse {
  expense_id: number;
  event_id: number;
  created_by: number;
  created_at: string;
  description: string;
  amount: number;
  currency: string;
  split_method: string;
  shares: ExpenseShare[];
}

export interface ExpensesResponse {
  items: ExpenseResponse[];
  total_count: number;
}

export interface UserBalance {
  user_id: number;
  username: string;
  balance: number;
  paid_amount: number;
  unpaid_amount: number;
  total_due: number;
}

export interface BalanceReportResponse {
  event_id: number;
  total_amount: number;
  user_balances: UserBalance[];
}