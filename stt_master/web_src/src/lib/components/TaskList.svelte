<script lang="ts">
    import type { ITask } from "../types";
    import TaskItem from "./TaskItem.svelte";

    interface Props {
        tasks: readonly ITask[];
    }

    let { tasks }: Props = $props()
</script>

<div class="task-list">
    {#each tasks as task (task.id) }
        <TaskItem {task} />
    {/each}
</div>

<style>
    .task-list {
        max-width: 700px;
        margin: 0 auto;
    }
</style>