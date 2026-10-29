<script lang="ts">
    import type { ITask } from "../types";
    import { TASK_STATUS } from "../types";
    import { formateDate } from "../utils";
    import StatusLabel from "./StatusLabel.svelte";

    interface Props  {
        task: ITask;
    }

    let { task }: Props = $props();
</script>

<div class="task-item">
    <div>ID: {task.id} ---- Время изменения: {formateDate(task.updated_at)} --- <StatusLabel status={task.status}/></div>
    <div>Файл: {task.file_path}</div>
</div>

<style>
    .task-item {
        padding: 1rem;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        margin-bottom: 0.75rem;
        background: white;
        transition: box-shadow 0.2s;
    }
</style>