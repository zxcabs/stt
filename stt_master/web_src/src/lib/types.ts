export const TASK_STATUS = {
    PENDING: 'pending',
    PROCESSING: 'processing',
    COMPLETED: 'completed',
    FAILED: 'failed'
} as const;

export type TaskStatus = typeof TASK_STATUS[keyof typeof TASK_STATUS];


export interface ITask {
    readonly id: number;
    readonly file_path: string;
    readonly status: TaskStatus;
    readonly error_message?: string;
    readonly updated_at: string;
}

export type ITaskStatusResponse = readonly ITask[];