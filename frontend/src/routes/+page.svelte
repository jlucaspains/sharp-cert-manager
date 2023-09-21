<script>
	import { onMount } from 'svelte';
	import { crossfade, fade, scale } from 'svelte/transition';
	import { PUBLIC_API_BASE_PATH } from '$env/static/public';
	import ValidationIcon from './validationIcon.svelte';

	class Item {
		url = '';
		name = '';
		issuer = '';
		signature = '';
		isValid = false;
		validity = 0;
		validityWarning = false;
		isLoaded = false;
		commonName = '';
		certStartDate = new Date();
		certEndDate = new Date();
		isCA = false;
		/**
		 * @type {string[]}
		 */
		certDnsNames = [];
		/**
		 * @type {{commonName: string, issuer: string, isCA: boolean}[]}
		 */
		otherCerts = [];
		/**
		 * @type {string[]}
		 */
		validationIssues = [];
	}

	/**
	 * @type {Item[]}
	 */
	let items = [];

	/**
	 * @type {Item | null}
	 */
	let selected = null;
	let showModal = false;

	const [send, receive] = crossfade({
		duration: 200,
		// @ts-ignore
		fallback: scale
	});

	onMount(async () => {
		const res = await fetch(`${PUBLIC_API_BASE_PATH}/site-list`);
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
				validity: 0,
				validityWarning: false,
				isLoaded: false,
				commonName: '',
				isCA: false,
				certDnsNames: [],
				otherCerts: [],
				validationIssues: [],
				certStartDate: new Date(),
				certEndDate: new Date()
			};
		});

		for (const item of items) {
			loadItem(item);
		}
	});

	/**
	 * @param {Item} item
	 */
	function loadItem(item) {
		fetch(`${PUBLIC_API_BASE_PATH}/check-url?url=${item.url}`)
			.then((res) => res.json())
			.then((data) => {
				item.issuer = data.issuer;
				item.signature = data.signature;
				item.isValid = data.isValid;
				item.commonName = data.commonName;
				item.otherCerts = data.otherCerts;
				item.certDnsNames = data.certDnsNames;
				item.isCA = data.isCA;
				item.certStartDate = new Date(data.certStartDate);
				item.certEndDate = new Date(data.certEndDate);
				item.validationIssues = data.validationIssues;
				item.validityWarning = data.expirationWarning;

				const validity = item.certEndDate.getTime() - new Date().getTime();
				item.validity = validity > 0 ? Math.floor(validity / (1000 * 3600 * 24)) : 0;

				item.isLoaded = true;
				items = items; // for display refresh
			});
	}

	/**
	 * @param {Item | null} item
	 */
	function select(item) {
		selected = item;
		showModal = true;
	}

	/**
	 * @param {Item | null} item
	 */
	function getItemState(item) {
		if (!item) return 'failed';

		if (!item.isLoaded) return 'loading';
		if (item.isValid) {
			if (item.validityWarning) return 'warning';
			else return 'success';
		}
		return 'failed';
	}
</script>

<div
	class="grid grid-flow-row gap-8 mt-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
>
	{#each items as item}
		<div
			class="check-item rounded shadow-lg shadow-gray-200 dark:shadow-gray-900 bg-white dark:bg-gray-800 duration-300 hover:-translate-y-1"
			in:receive={{ key: item.name }}
			out:send={{ key: item.name }}
			on:click={() => select(item)}
			on:keydown={({ key }) => {
				if (key === 'enter') select(item);
			}}
			role="button"
			tabindex="0"
		>
			<div class="flex rounded-lg h-full bg-gray-600 p-8 flex-col">
				<div class="flex items-center mb-3">
					<ValidationIcon type={getItemState(item)} />
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
	{#if showModal}
		<div
			id="defaultModal"
			tabindex="-1"
			class="fixed top-0 left-0 right-0 z-50 w-full p-4 overflow-x-hidden overflow-y-auto md:inset-0 h-[calc(100%-1rem)] max-h-full justify-center items-center flex"
			in:receive|global={{ key: selected?.name }}
			out:send|global={{ key: selected?.name }}
		>
			<div class="relative w-full max-w-2xl max-h-full">
				<!-- Modal content -->
				<div class="relative bg-white rounded-lg shadow dark:bg-gray-700">
					<!-- Modal header -->
					<div class="flex items-start justify-between p-4 border-b rounded-t dark:border-gray-600">
						<ValidationIcon type={selected?.isValid ? 'success' : 'failed'} />
						<h3 class="text-xl font-semibold text-gray-900 dark:text-white">
							{selected?.name}
						</h3>
						<button
							type="button"
							class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ml-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white"
							on:click={() => {
								showModal = false;
							}}
						>
							<svg
								class="w-3 h-3"
								aria-hidden="true"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 14 14"
							>
								<path
									stroke="currentColor"
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
								/>
							</svg>
							<span class="sr-only">Close modal</span>
						</button>
					</div>
					<!-- Modal body -->
					<div class="p-6 space-y-6">
						<table class="table-auto leading-relaxed text-base text-gray-400">
							<tbody>
								<tr class="border-gray-700">
									<td class="px-4 py-2 text-white">Common Name</td>
									<td class="px-4 py-2">{selected?.commonName}</td>
								</tr>
								<tr>
									<td class="px-4 py-2 text-white">Issuer</td>
									<td class="px-4 py-2">{selected?.issuer}</td>
								</tr>
								<tr>
									<td class="px-4 py-2 text-white">Signature</td>
									<td class="px-4 py-2">{selected?.signature}</td>
								</tr>
								<tr>
									<td class="px-4 py-2 text-white">Validity</td>
									<td class="px-4 py-2">
										<ul class="">
											<li>{selected?.validity} days</li>
											<li>Issued on: {selected?.certStartDate.toLocaleString()}</li>
											<li>Expires on: {selected?.certEndDate.toLocaleString()}</li>
										</ul></td
									>
								</tr>
								<tr>
									<td class="px-4 py-2 text-white">Is CA</td>
									<td class="px-4 py-2">{selected?.isCA}</td>
								</tr>
								<tr>
									<td class="px-4 py-2 text-white">DNS Names</td>
									<td class="px-4 py-2">{selected?.certDnsNames.join(', ')}</td>
								</tr>
								{#each selected?.otherCerts || [] as cert, i}
									<tr>
										<td class="px-4 py-2 text-white">Other cert {i + 1}</td>
										<td class="px-4 py-2"
											><ul class="">
												<li>Common Name: <span>{cert.commonName}</span></li>
												<li>Issuer: {cert.issuer}</li>
												<li>Is CA: {cert.isCA}</li>
											</ul></td
										>
									</tr>
								{/each}
								<tr>
									<td class="px-4 py-2 text-white">Validation</td>
									<td class="px-4 py-2">
										{#if selected?.isValid}
											Cert is valid
										{:else}
											<ul class="">
												{#each selected?.validationIssues || [] as issue}
													<li>{issue}</li>
												{/each}
											</ul>
										{/if}
									</td>
								</tr>
							</tbody>
						</table>
					</div>
					<!-- Modal footer -->
					<div
						class="flex items-center p-6 space-x-2 border-t border-gray-200 rounded-b dark:border-gray-600"
					>
						<button
							type="button"
							on:click={() => {
								showModal = false;
							}}
							class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
							>OK</button
						>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	button,
	.check-item {
		will-change: transform;
	}
</style>
