<script lang="ts">
	import './layout.css'
	import favicon from '$lib/assets/favicon.svg'
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query'
	import { browser } from '$app/environment'

	let { children } = $props()

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser, // Only run in browser (SSR safe)
				staleTime: 1000 * 60, // 1 minute
				refetchOnWindowFocus: false
			}
		}
	})
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<QueryClientProvider client={queryClient}>
	{@render children()}
</QueryClientProvider>
