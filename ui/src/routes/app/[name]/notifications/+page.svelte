<script>
  import { page } from "$app/stores";
  import DiscordWebhook from "$components/notifications/DiscordWebhook.svelte";
</script>

<DiscordWebhook app={$page.params.name} />
