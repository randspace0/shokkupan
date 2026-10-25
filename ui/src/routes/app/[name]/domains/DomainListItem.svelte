<script>
  import { page } from "$app/stores";
  import { removeAppDomain } from "$lib/api";
  import { createEventDispatcher } from "svelte";

  import Icon from "$common/Icon.svelte";
  import ConfirmationModal from "$common/ConfirmationModal.svelte";
  import Alert from "$common/Alert.svelte";

  export let domain = "";
  export let loading = false;

  let domainHref = domain.startsWith("http") ? domain : `https://${domain}`;
  const dispatch = createEventDispatcher();
  const signalRemoved = () => dispatch("removeDomain", { domain });

  let confirmationModalOpen = false;
</script>

<li class="list-row items-center">
  <a class="link link-hover list-col-grow text-lg" href={domainHref}>{domain}</a>
  <button
    class="btn btn-ghost btn-sm btn-circle text-error"
    aria-label="Remove {domain}"
    disabled={loading}
    on:click={() => (confirmationModalOpen = true)}
  >
    {#if loading}
      <span class="loading loading-spinner loading-xs"></span>
    {:else}
      <Icon type="delete" size="sm" />
    {/if}
  </button>
</li>

<ConfirmationModal
  name="remove-domain-modal"
  title="Removing Domain"
  action="remove domain {domain}"
  on:accepted={signalRemoved}
  bind:open={confirmationModalOpen}
  doingAction={loading}
/>
