<script lang="ts">
    import type { TaskStatus } from "../types";
    import { TASK_STATUS } from "../types";

    interface Props {
        status: TaskStatus;
    }

    let { status }: Props = $props();

    const labels: Record<TaskStatus, string>= {
        [TASK_STATUS.PENDING]: "В ожидании",
        [TASK_STATUS.PROCESSING]: "В процессе",
        [TASK_STATUS.COMPLETED]: "Завершено",
        [TASK_STATUS.FAILED]: "Ошибка"
    }
</script>

<span class="status-label {status}">{labels[status]}</span>

<style>
    .status-label {
        display: inline-block;
        padding: 0.25rem 0.6rem;
        border-radius: 12px;
        font-size: 0.8rem;
        font-weight: 500;
        border: 1px solid;

    }
    
    .status-label.pending {
        background-color: #f3f4f6;
        color: #6b7280;
        border-color: #d1d5db;
    }

    .status-label.processing {
        background-color: #fef3c7;
        color: #d97706;
        border-color: #fbbf24;
    }

    .status-label.completed {
        background-color: #d1fae5;
        color: #059669;
        border-color: #6ee7b7;
    }

    .status-label.failed {
        background-color: #fee2e2;
        color: #dc2626;
        border-color: #fca5a5;
    }
</style>