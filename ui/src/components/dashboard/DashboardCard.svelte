<script>
  import Icon from "$common/Icon.svelte";

  export let info;
  export let contentType;

  let state;
  let badge;
  $: if (info) {
    if (info["is_running"]) [state, badge] = ["Running", "badge-success"];
    else if (info["is_deployed"]) [state, badge] = ["Deployed", "badge-warning"];
    else if (!info["is_setup"]) [state, badge] = ["Needs setup", "badge-info"];
    else [state, badge] = ["Stopped", "badge-ghost"];
  }
</script>

<a
  href="/{contentType}/{info.name}"
  class="card card-border bg-base-100 rounded-3xl hover:border-primary/60 hover:shadow-lg hover:-translate-y-0.5 transition motion-reduce:transition-none motion-reduce:hover:translate-y-0 focus-visible:outline-2 focus-visible:outline-primary"
>
  <div class="card-body flex-row items-center gap-3 p-4">
    {#if info.type}
      <Icon size="lg" type={info.type} />
    {/if}
    <span class="font-semibold text-lg truncate flex-1">{info.name}</span>
    {#if contentType === "app" && state}
      <span class="badge badge-soft {badge} gap-1">
        <span class="status {badge.replace('badge-', 'status-')}"></span>
        {state}
      </span>
    {/if}
  </div>
</a>
