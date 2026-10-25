<script>
  import { goto } from "$app/navigation";

  import Icon from "$common/Icon.svelte";
  import QueryDataWrapper from "$common/QueryDataWrapper.svelte";
  import DashboardCard from "./DashboardCard.svelte";

  export let contentType = "";
  export let states = {};
  export let query;

  let contentTypeTitle;
  $: contentTypeTitle =
    contentType.substring(0, 1).toUpperCase() + contentType.substring(1);
</script>

<QueryDataWrapper {query} action="loading {contentTypeTitle}s">
  <section class="flex flex-col gap-3" aria-label="{contentTypeTitle}s">
    <header class="flex items-center gap-3">
      <h2 class="text-xl font-bold flex-1">
        {contentTypeTitle}s
        <span class="badge badge-neutral badge-sm align-middle ml-1"
          >{$query.data.length}</span
        >
      </h2>
      <button
        class="btn btn-primary btn-sm gap-2 min-h-11 sm:min-h-8"
        on:click={() => goto(`/${contentType}/new`)}
      >
        <Icon size="sm" type="plus" />
        Create
      </button>
    </header>

    {#if $query.data.length > 0}
      <div class="flex flex-col gap-2">
        {#each $query.data as info}
          <DashboardCard {contentType} {info} />
        {/each}
      </div>
    {:else}
      <div class="card card-border border-dashed bg-base-100">
        <div class="card-body items-center text-center">
          <p class="text-base-content/70">No {contentType}s yet.</p>
          <button
            class="btn btn-primary gap-2"
            on:click={() => goto(`/${contentType}/new`)}
          >
            <Icon size="md" type="add" />
            Create {contentTypeTitle}
          </button>
        </div>
      </div>
    {/if}
  </section>
</QueryDataWrapper>
