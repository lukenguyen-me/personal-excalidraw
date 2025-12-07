<script lang="ts">
	import { authStore } from '$lib/stores/auth'

	let { onSuccess } = $props<{ onSuccess: () => void }>()

	let modal = $state<HTMLDialogElement>()
	let keyInput = $state('')
	let error = $state('')
	let loading = $state(false)

	// Auto-open modal when mounted
	$effect(() => {
		if (modal) {
			modal.showModal()
		}
	})

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault()

		if (!keyInput.trim()) {
			error = 'Please enter an access key'
			return
		}

		loading = true
		error = ''

		try {
			// Validate key by calling auth endpoint
			const response = await fetch('/api/auth/validate', {
				headers: {
					Authorization: `Bearer ${keyInput}`
				}
			})

			if (response.ok) {
				authStore.setAccessKey(keyInput)
				modal?.close()
				onSuccess()
			} else {
				const data = await response.json()
				error = data.message || 'Invalid access key'
			}
		} catch (err) {
			error = 'Network error. Please try again.'
		} finally {
			loading = false
		}
	}
</script>

<dialog bind:this={modal} class="modal" oncancel={(e) => e.preventDefault()}>
	<div class="modal-box">
		<h3 class="text-lg font-bold mb-4">Access Key Required</h3>
		<p class="mb-4">Please enter your access key to continue.</p>

		<form onsubmit={handleSubmit}>
			<input
				type="password"
				bind:value={keyInput}
				class="input input-bordered w-full mb-2"
				placeholder="Enter access key"
				disabled={loading}
			/>

			{#if error}
				<div class="alert alert-error alert-soft mb-4">
					<span>{error}</span>
				</div>
			{/if}

			<div class="modal-action">
				<button type="submit" class="btn btn-primary" disabled={loading}>
					{loading ? 'Validating...' : 'Submit'}
					{#if loading}
						<span class="loading loading-spinner loading-sm"></span>
					{/if}
				</button>
			</div>
		</form>
	</div>
</dialog>
