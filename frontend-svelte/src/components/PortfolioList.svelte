<script>
  import { onMount } from 'svelte';
  import { fetchPortfolio } from '../lib/api.js';

  let projects = [];
  let loading = true;
  let error = null;

  onMount(async () => {
    try {
      projects = await fetchPortfolio();
    } catch (e) {
      error = 'Failed to load projects.';
    } finally {
      loading = false;
    }
  });
</script>

{#if loading}
  <p>Loading portfolio...</p>
{:else if error}
  <p>{error}</p>
{:else}
  <ul>
    {#each projects as project}
      <li>
        <h3>{project.name}</h3>
        <p>{project.description}</p>
      </li>
    {/each}
  </ul>
{/if}

<style>
  ul {
    list-style: none;
    padding: 0;
  }
  li {
    margin-bottom: 1rem;
    border-bottom: 1px solid #ccc;
    padding-bottom: 0.5rem;
  }
</style>
