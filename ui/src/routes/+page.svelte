<script>
  import { useQuery } from "$lib/query";
  import { getAllAppsOverview, listServices } from "$lib/api";

  import DashboardCardList from "$components/dashboard/DashboardCardList.svelte";

  const queryOpts = {
    refetchOnMount: true,
  };
  const appsOverview = useQuery(
    "getAllAppsOverview",
    getAllAppsOverview,
    queryOpts
  );
  const servicesList = useQuery("listServices", listServices, queryOpts);
  const appStates = {};
  const serviceStates = {};
</script>

<div class="h-full overflow-y-auto">
  <div class="mx-auto max-w-6xl p-4 grid gap-8 md:grid-cols-2 items-start">
    <DashboardCardList contentType="app" query={appsOverview} states={appStates} />
    <DashboardCardList
      contentType="service"
      query={servicesList}
      states={serviceStates}
    />
  </div>
</div>
