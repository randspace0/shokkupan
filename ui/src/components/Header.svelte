<script>
  import Icon from "$common/Icon.svelte";
  import { createEventDispatcher } from "svelte";

  export let isAuthenticated = false;
  export let username = "";
  // "" until known; Settings is admin-only
  export let role = "";
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
  class="bg-base-100/80 backdrop-blur text-base-content border border-base-300 shadow-sm rounded-3xl flex flex-wrap items-center gap-2 px-4 py-2 mx-3 mt-3"
>
  <a
    class="inline-flex items-center justify-center gap-2 btn btn-ghost normal-case text-xl min-h-11 focus-visible:outline-2 focus-visible:outline-offset-2"
    href="/"
    ><img src="/favicon.svg" alt="" class="size-7" />shokkupan</a
  >

  {#if isAuthenticated}
    <ul class="flex flex-wrap items-center gap-1 text-sm font-medium">
      {#each authenticatedLinks.filter((l) => l.href !== "/logout" && (l.href !== "/settings" || role === "admin")) as link}
        <li><a class={linkStyle} href={link["href"]}>{link["text"]}</a></li>
      {/each}
    </ul>
  {/if}

  <div class="flex-1"></div>

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

  {#if isAuthenticated}
    <div class="flex items-center gap-2" title={username ? `Signed in as ${username}` : ""}>
      {#if username}
        <div class="avatar avatar-placeholder">
          <div class="bg-primary text-primary-content w-8 rounded-full">
            <span class="text-xs uppercase">{username.slice(0, 2)}</span>
          </div>
        </div>
        <span class="font-semibold hidden sm:inline">{username}</span>
        {#if role}<span class="badge badge-soft badge-sm hidden sm:inline-flex">{role}</span>{/if}
      {/if}
      <a class={linkStyle} href="/logout">Logout</a>
    </div>
  {/if}
</nav>
