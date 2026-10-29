<script lang="ts">
  import { fetchTaskStatus } from "./lib/api";
  import type { ITask } from "./lib/types";
  import TaskList from "./lib/components/TaskList.svelte";
    import { onMount } from "svelte";

  let tasks: readonly ITask[] = $state([]);
  let loading: boolean = $state(false);
  let error: string | null = $state(null);

  async function loadTasks() {
    if (loading) return;

    loading = true;
    error = null;

    try {
      tasks = await fetchTaskStatus();
    } catch (e) {
      error = e instanceof Error ? e.message : "Ошибка";
    } finally {
      loading = false
    }
  }

  onMount(() => {
    loadTasks();
  });
</script>

<div class="app">
  <header>
    <h1>Статус задач:</h1>
    <button onclick={loadTasks} disabled={loading}>{ loading ? "Загрузка...": "Обновить" }</button>
  </header>

  <main>
    {#if loading && tasks?.length === 0 }
      <div class="loading">Загрузка данных...</div>
    {:else if error} 
      <div class="error">
        <p>Ошибка: {error}</p>
        <button onclick={loadTasks}>Попробовать снова</button>
      </div>
    {:else}
      <TaskList {tasks}/>
    {/if}
  </main>
</div>

<style>
    :global(body) {
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        margin: 0;
        padding: 2rem;
        background: #f9fafb;
        color: #111827;
    }

    .app {
        max-width: 900px;
        margin: 0 auto;
    }

    header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 2rem;
    }

    h1 {
        margin: 0;
        font-size: 1.75rem;
    }

    button {
        padding: 0.5rem 1.25rem;
        background: #3b82f6;
        color: white;
        border: none;
        border-radius: 6px;
        cursor: pointer;
        font-size: 0.95rem;
        font-weight: 500;
        transition: background 0.2s;
    }

    button:hover:not(:disabled) {
        background: #2563eb;
    }

    button:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .loading, .error {
        text-align: center;
        padding: 2rem;
        background: white;
        border-radius: 8px;
        max-width: 700px;
        margin: 0 auto;
    }

    .error {
        color: #dc2626;
        border: 1px solid #fecaca;
    }

    .error button {
        margin-top: 1rem;
        background: #dc2626;
    }

    .error button:hover {
        background: #b91c1c;
    }
</style>