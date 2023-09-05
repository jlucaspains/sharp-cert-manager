<script>
	/**
	 * @type {'success' | 'failed' | 'loading' | 'warning'}
	 */
	export let type = 'success';
</script>

{#if type == 'loading'}
	<div
		class="w-8 h-8 mr-3 inline-flex items-center justify-center rounded-full bg-blue-600 text-white flex-shrink-0"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="w-6 h-6"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z"
			/>
		</svg>
	</div>
{:else if type == 'success'}
	<div
		class="w-8 h-8 mr-3 inline-flex items-center justify-center rounded-full bg-green-600 text-white flex-shrink-0"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="w-6 h-6"
		>
			<path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
		</svg>
	</div>
{:else if type == 'failed'}
	<div
		class="w-8 h-8 mr-3 inline-flex items-center justify-center rounded-full bg-red-600 text-white flex-shrink-0"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="w-6 h-6"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
			/>
		</svg>
	</div>
{:else}
	<div
		class="w-8 h-8 mr-3 inline-flex items-center justify-center rounded-full bg-yellow-600 text-white flex-shrink-0"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="w-6 h-6"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
			/>
		</svg>
	</div>
{/if}
