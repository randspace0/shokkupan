<script>
  import { page } from "$app/stores";
  import {
    useMutation,
    useQuery,
    useQueryClient,
  } from "$lib/query";
  import { addAppDomain, getAppDomainsReport, removeAppDomain } from "$lib/api";

  import QueryDataWrapper from "$common/QueryDataWrapper.svelte";
  import Card from "$common/Card.svelte";
  import Icon from "$common/Icon.svelte";

  import AddDomainModal from "./AddDomainModal.svelte";
  import DomainListItem from "./DomainListItem.svelte";

  const appName = $page.params.name;

  const queryClient = useQueryClient();
  const queryKey = [{ appName }, "getAppDomainsReport"];
  const domainsReport = useQuery(queryKey, () => getAppDomainsReport(appName));

  let addDomainModalOpen = false;

  const onSuccess = () => {
    addDomainModalOpen = false;
    queryClient.invalidateQueries(queryKey);
  };

  const addDomainMutation = useMutation(
    (domain) => addAppDomain(appName, domain),
    { onSuccess }
  );

  const removeDomainMutation = useMutation(
    (domain) => removeAppDomain(appName, domain),
    { onSuccess }
  );

  const addDomain = ({ detail }) => $addDomainMutation.mutate(detail.domain);
  const removeDomain = ({ detail }) =>
    $removeDomainMutation.mutate(detail.domain);

  let domains = [];
  $: if ($domainsReport.isSuccess && $domainsReport.data) {
    domains = $domainsReport.data["domains"] || [];
  }
</script>

<QueryDataWrapper query={domainsReport} action="fetching domains">
  <Card title="Domains">
    {#if domains.length === 0}
      <p class="text-base-content/70">No domains configured.</p>
    {:else}
      <ul class="list rounded-2xl border border-base-300 bg-base-100">
        {#each domains as domain}
          <DomainListItem {domain} on:removeDomain={removeDomain} />
        {/each}
      </ul>
    {/if}

    <div slot="actions">
      <button class="btn btn-primary gap-2" on:click={() => (addDomainModalOpen = true)}>
        Add Domain
        <Icon type="add" />
      </button>
    </div>
  </Card>
</QueryDataWrapper>

<AddDomainModal
  loading={$addDomainMutation.isLoading}
  bind:open={addDomainModalOpen}
  on:addDomain={addDomain}
/>
