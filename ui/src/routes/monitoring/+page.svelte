<script>
  import { getMonitoring, getServerMonitoring } from "$lib/api";
  import { useQuery } from "$lib/query";
  import Card from "$common/Card.svelte";
  import QueryDataWrapper from "$common/QueryDataWrapper.svelte";
  import AppStats from "$components/monitoring/AppStats.svelte";
  import ServerStats from "$components/monitoring/ServerStats.svelte";

  const monitoring = useQuery("getMonitoring", getMonitoring, {
    refetchInterval: 30000,
  });

  const serverMonitoring = useQuery("getServerMonitoring", getServerMonitoring);
</script>

<div class="flex flex-col gap-4 p-4">
  {#if $serverMonitoring.data}
    <ServerStats />
  {/if}
  <QueryDataWrapper query={monitoring} action="loading monitoring data">
    <Card title="Apps">
      <div class="flex flex-col gap-4">
        {#each $monitoring.data.apps as app}
          <div class="rounded-2xl border border-base-300 p-4 flex flex-col gap-3">
            <div class="flex flex-wrap items-center gap-3">
              <a class="link link-hover font-semibold text-lg" href={`/app/${app.name}`}>{app.name}</a>
              {#if !app.is_deployed}
                <span class="badge badge-ghost">not deployed</span>
              {:else if app.is_running}
                <span class="badge badge-soft badge-success">running</span>
              {:else}
                <span class="badge badge-soft badge-error">stopped</span>
              {/if}
              <span class="text-sm text-base-content/70">
                {app.num_processes} process{app.num_processes === 1 ? "" : "es"}
              </span>
            </div>
            {#if app.is_running && $monitoring.data.stats_available}
              <AppStats app={app.name} />
            {/if}
          </div>
        {/each}

        {#if !$monitoring.data.stats_available}
          <div class="alert alert-info alert-soft items-start">
            <div>
              <p>Install the stats plugin to chart CPU and memory per app. On the server run:</p>
              <code class="block bg-base-300 p-2 rounded-xl mt-2 text-xs break-all">
                curl -fsSL https://raw.githubusercontent.com/randspace0/shokkupan/main/plugins/install.sh | sudo bash -s stats
              </code>
              <p class="mt-2">Samples appear after the next poll (up to a minute).</p>
            </div>
          </div>
        {/if}
      </div>
    </Card>

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
