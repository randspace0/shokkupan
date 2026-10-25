<script>
  import { page } from "$app/stores";

  export let prefix = "";
  export let pages = [];

  const prefixLen = prefix.length + 1;
  const getCurrentPage = () => $page.url.pathname.slice(prefixLen);
  let currentPath = "";
  let currentPage = getCurrentPage();
  let pageLabel = "";
  let openMobileMenu;
  $: if ($page.url.pathname !== currentPath) {
    currentPath = $page.url.pathname;
    currentPage = getCurrentPage();
    for (let i = 0; i < pages.length; i++) {
      let info = pages[i];
      if (info.path === currentPage) pageLabel = info.name;
    }
    openMobileMenu = false;
  }
</script>

<nav
  aria-label="Section"
  class="menu rounded-box hidden sm:block w-48 shrink-0 max-h-full bg-base-100 border border-base-300 overflow-y-auto h-fit"
>
  {#each pages as info, i}
    <li class="">
      <a href="{prefix}/{info.path}" class:menu-active={info.path === currentPage}>
        <span>{info.name}</span>
      </a>
    </li>
  {/each}
</nav>

<div class="inline-block sm:hidden w-full">
  <div class="collapse collapse-arrow bg-base-100 border border-base-300 w-full">
    <input type="checkbox" class="peer" bind:checked={openMobileMenu} />
    <div class="collapse-title text-xl font-medium peer-checked:bg-base-300">
      {pageLabel}
    </div>
    <div class="collapse-content peer-checked:bg-base-300">
      <div class="menu w-full">
        {#each pages as info, i}
          <li>
            <a
              href="{prefix}/{info.path}"
              class:menu-active={info.path === currentPage}
            >
              <span>{info.name}</span>
            </a>
          </li>
        {/each}
      </div>
    </div>
  </div>
</div>
