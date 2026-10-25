<script>
  import { createEventDispatcher } from "svelte";

  import Icon from "$common/Icon.svelte";
  import ConfirmationModal from "$common/ConfirmationModal.svelte";

  import LinkModal from "./LinkModal.svelte";

  export let serviceName;
  export let appName;
  export let serviceType;
  export let cardType = "app";
  export let isLinked = false;
  export let loading;

  const dispatch = createEventDispatcher();

  let cardLink;
  $: cardLink =
    cardType === "service" ? `/service/${serviceName}` : `/app/${appName}`;

  let linkModalOpen = false;
  const doLink = ({ detail }) => {
    dispatch("link", { serviceName, appName, options: detail });
  };

  let unlinkModalOpen = false;
  const doUnlink = () => {
    dispatch("unlink", { serviceName, appName });
  };
</script>

<div
  class="flex flex-row w-full sm:w-auto p-3 bg-base-200 border border-base-300 items-center gap-3 rounded-2xl"
>
  <a class="btn btn-ghost" class:gap-4={cardType === "service"} href={cardLink}>
    <div class="">
      {#if cardType === "service"}<Icon type={serviceType} size="lg" />{/if}
    </div>
    <div class="">
      <span>{cardType === "service" ? serviceName : appName}</span>
    </div>
  </a>
  <div class="flex-grow" />
  <div class:hidden={!isLinked}>
    <button
      class="btn btn-outline btn-error gap-2"
      disabled={loading}
      on:click={() => (unlinkModalOpen = true)}
    >
      {#if loading}
        <span class="loading loading-spinner loading-xs"></span>
        Unlinking
      {:else}
        <Icon type="unlink" />
        Unlink
      {/if}
    </button>
  </div>
  <div class:hidden={isLinked}>
    <button
      class="btn btn-primary gap-2"
      disabled={loading}
      on:click={() => (linkModalOpen = true)}
    >
      {#if loading}
        <span class="loading loading-spinner loading-xs"></span>
        Linking
      {:else}
        <Icon type="link" size="lg" />
        Link
      {/if}
    </button>
  </div>
</div>

{#if isLinked}
  <ConfirmationModal
    name="unlink"
    title="Confirm Unlinking"
    action="unlink service '{serviceName}' and app '{appName}'"
    bind:open={unlinkModalOpen}
    on:accepted={doUnlink}
  />
{:else}
  <LinkModal
    {serviceType}
    preventClose={loading}
    bind:open={linkModalOpen}
    on:link={doLink}
  />
{/if}
