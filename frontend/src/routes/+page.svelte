<script>
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';

	/**
	 * @type {{url: string, name: string, issuer: string, signature: string, isValid: boolean, validationError: string, validity: number, isLoaded: boolean}[]}
	 */
	let items = [];

	onMount(async () => {
		const res = await fetch(`http://localhost:3000/api/site-list`);
		/**
		 * @type {{url: string, name: string}[]}
		 */
		const item = await res.json();

		items = item.map((item) => {
			return {
				...item,
				issuer: '',
				signature: '',
				isValid: false,
				validationError: '',
				validity: 0,
				isLoaded: false
			};
		});

		for (const item of items) {
			loadItem(item);
		}
	});

	/**
	 * @param {{ url: any; issuer: any; signature: any; isValid: any; validationError: any; validity: any; isLoaded: boolean; }} item
	 */
	function loadItem(item) {
		fetch(`http://localhost:3000/api/check-url?url=${item.url}`)
			.then((res) => res.json())
			.then((data) => {
				console.log(data);
				item.issuer = data.issuer;
				item.signature = data.signature;
				item.isValid = data.isValid;
				item.validationError = data.validationError;

				const startDate = new Date(data.certStartDate);
				const endDate = new Date(data.certEndDate);
				const validity = endDate.getTime() - startDate.getTime();
				item.validity = validity > 0 ? Math.floor(validity / (1000 * 3600 * 24)) : 0;

				item.isLoaded = true;
				items = items; // for display refresh
			});
	}
</script>

<div
	class="grid grid-flow-row gap-8 mt-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
>
	{#each items as item}
		<div
			class="rounded shadow-lg shadow-gray-200 dark:shadow-gray-900 bg-white dark:bg-gray-800 duration-300 hover:-translate-y-1"
		>
			<div class="flex rounded-lg h-full bg-gray-600 p-8 flex-col">
				<div class="flex items-center mb-3">
					{#if !item.isLoaded}
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
					{:else if item.isValid}
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
					{:else}
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
					{/if}
					<h2 class="text-white text-lg font-medium">{item.name}</h2>
				</div>
				{#if item.isLoaded}
					<div class="flex flex-col justify-between flex-grow">
						<p transition:fade class="leading-relaxed text-base text-white">
							Issuer: <span class="font-bold">{item.issuer}</span>
						</p>
						<p transition:fade class="leading-relaxed text-base text-white">
							Signature: <span class="font-bold">{item.signature}</span>
						</p>
						<p transition:fade class="leading-relaxed text-base text-white">
							Validity: <span class="font-bold">{item.validity} days</span>
						</p>
						<a
							transition:fade
							href="#/view"
							class="mt-3 text-white hover:text-blue-600 inline-flex items-center"
							>Review
							<svg
								fill="none"
								stroke="currentColor"
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								class="w-4 h-4 ml-2"
								viewBox="0 0 24 24"
							>
								<path d="M5 12h14M12 5l7 7-7 7" />
							</svg>
						</a>
					</div>
				{:else}
					<div role="status" class="max-w-sm animate-pulse">
						<div class="h-2.5 bg-gray-200 rounded-full dark:bg-gray-700 w-48 mb-4" />
						<div class="h-2 bg-gray-200 rounded-full dark:bg-gray-700 max-w-[360px] mb-2.5" />
						<div class="h-2 bg-gray-200 rounded-full dark:bg-gray-700 mb-2.5" />
						<div class="h-2 bg-gray-200 rounded-full dark:bg-gray-700 max-w-[330px] mb-2.5" />
						<div class="h-2 bg-gray-200 rounded-full dark:bg-gray-700 max-w-[300px] mb-2.5" />
						<div class="h-2 bg-gray-200 rounded-full dark:bg-gray-700 max-w-[360px]" />
						<span class="sr-only">Loading...</span>
					</div>
				{/if}
			</div>
		</div>
	{/each}
</div>
