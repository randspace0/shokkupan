<script>
  import { getMonitoring } from "$lib/api";
  import { useQuery } from "@sveltestack/svelte-query";
  import Card from "$common/Card.svelte";
  import QueryDataWrapper from "$common/QueryDataWrapper.svelte";
  import AppStats from "$components/monitoring/AppStats.svelte";

  const monitoring = useQuery("getMonitoring", getMonitoring, {
    refetchInterval: 30000,
  });
</script>

<div class="flex flex-col gap-4 p-4">
  <QueryDataWrapper query={monitoring} action="loading monitoring data">
    <Card title="Apps">
      <div class="overflow-x-auto">
        <table class="table w-full">
          <thead>
            <tr><th>App</th><th>Status</th><th>Processes</th></tr>
          </thead>
          <tbody>
            {#each $monitoring.data.apps as app}
              <tr>
                <td><a class="link" href={`/app/${app.name}`}>{app.name}</a></td>
                <td>
                  {#if !app.is_deployed}
                    <span class="badge">not deployed</span>
                  {:else if app.is_running}
                    <span class="badge badge-success">running</span>
                  {:else}
                    <span class="badge badge-error">stopped</span>
                  {/if}
                </td>
                <td>{app.num_processes}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </Card>

    {#if $monitoring.data.stats_available}
      {#each $monitoring.data.apps.filter((a) => a.is_running) as app}
        <AppStats app={app.name} />
      {/each}
    {:else}
      <Card title="CPU and memory">
        <p>
          Install the dokku-stats plugin to chart CPU and memory. On the server run:
        </p>
        <code class="block bg-base-300 p-2 rounded">
          curl -fsSL https://raw.githubusercontent.com/randspace0/shokkupan/main/plugins/install.sh | sudo bash -s stats
        </code>
        <p>Samples appear after the next poll (up to a minute).</p>
      </Card>
    {/if}

    <Card title="Recent status changes (30 days)">
      {#if $monitoring.data.events.length === 0}
        <p>No changes recorded yet.</p>
      {:else}
        <ul>
          {#each $monitoring.data.events as ev}
            <li>
              {new Date(ev.at).toLocaleString()}: <b>{ev.app_name}</b>
              {ev.running ? "started running" : "stopped"}
            </li>
          {/each}
        </ul>
      {/if}
    </Card>
  </QueryDataWrapper>
</div>
