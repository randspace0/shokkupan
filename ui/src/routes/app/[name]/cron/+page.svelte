<script>
  import { page } from "$app/stores";
  import { getAppCron, runAppCron } from "$lib/api";
  import { commandExecutionIds, executionIdDescriptions } from "$lib/stores";
  import { useMutation, useQuery } from "@sveltestack/svelte-query";
  import Card from "$common/Card.svelte";
  import Error from "$common/Error.svelte";
  import QueryDataWrapper from "$common/QueryDataWrapper.svelte";

  const appName = $page.params.name;
  const tasks = useQuery([{ appName }, "getAppCron"], () => getAppCron(appName));

  const run = useMutation(async (task) => {
    const id = await runAppCron(appName, task.id);
    $executionIdDescriptions[id] = `Run cron task ${task.command}`;
    return commandExecutionIds.addID(id);
  });
</script>

<QueryDataWrapper query={tasks} action="loading cron tasks">
  <Card title="Scheduled jobs">
    {#if $tasks.data.length === 0}
      <p>
        No cron tasks. Define them in the <code>cron</code> section of your app's
        <code>app.json</code> and redeploy.
      </p>
    {:else}
      <div class="overflow-x-auto">
        <table class="table w-full">
          <thead>
            <tr><th>Schedule</th><th>Command</th><th /></tr>
          </thead>
          <tbody>
            {#each $tasks.data as task}
              <tr>
                <td><code>{task.schedule}</code></td>
                <td><code>{task.command}</code></td>
                <td>
                  <button
                    class="btn btn-sm"
                    disabled={$run.isLoading}
                    on:click={() => $run.mutate(task)}>Run now</button
                  >
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
    {#if $run.isError}<Error action="running cron task" error={$run.error} />{/if}
  </Card>
</QueryDataWrapper>
