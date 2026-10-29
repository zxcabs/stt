import type { ITaskStatusResponse } from "./types";

export class ApiError extends Error {
    name: string = "ApiError";

    constructor(public status: number, message: string) {
        super(message);
    } 
}

export async function fetchTaskStatus(): Promise<ITaskStatusResponse> {
    const response = await fetch("/api/tasks/status");

    if (!response.ok) {
        throw new ApiError(response.status, "Ошибка загрузки")
    }

    return await response.json() as ITaskStatusResponse
}