<script>
  import Icon from "$common/Icon.svelte";
  import { createEventDispatcher } from "svelte";

  export let isAuthenticated = false;
  export let showExecutionButton = false;
  export let executionWindowVisible = false;

  const authenticatedLinks = [
    { href: "/monitoring", text: "Monitoring" },
    { href: "/settings", text: "Settings" },
    { href: "/logout", text: "Logout" },
  ];
  const linkStyle = "px-3 py-2 btn btn-ghost min-h-11 focus-visible:outline-2 focus-visible:outline-offset-2";

  const dispatch = createEventDispatcher();
  const executionButtonClicked = () => dispatch("executionButtonClicked");
</script>

<nav
  aria-label="Main"
  class="bg-neutral text-neutral-content flex flex-col md:flex-row items-center w-full justify-between gap-2 p-4 mx-auto"
>
  <a
    class="inline-flex items-center justify-center btn btn-ghost normal-case text-xl min-h-11 focus-visible:outline-2 focus-visible:outline-offset-2"
    href="/">shokkupan</a
  >

  {#if showExecutionButton}
    <button
            class="btn gap-2 min-h-11 hover:btn-secondary focus-visible:outline-2 focus-visible:outline-offset-2"
            class:btn-active={executionWindowVisible}
            on:click={executionButtonClicked}
    >
      <Icon type="file-text" />
      command output
    </button>
  {/if}

  <div class="flex flex-col">
    <ul class="flex flex-col md:flex-row items-center gap-2 text-sm font-medium">
      {#if isAuthenticated}
        {#each authenticatedLinks as link}
          <li>
            <a class={linkStyle} href={link["href"]}>{link["text"]}</a>
          </li>
        {/each}
      {/if}
    </ul>
  </div>
</nav>
