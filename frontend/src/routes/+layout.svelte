<script>
	// import Header from './Header.svelte';
	import './styles.css';
</script>

<div class="app bg-gray-900">
	<header class="w-full bg-gray-900 shadow-sm body-font">
		<div class="flex flex-col flex-wrap items-center p-5 mx-auto">
			<div class="top-bar-title mb-4 text-white">sharp-cert-manager</div>
		</div>
	</header>

	<main>
		<slot />
	</main>
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	main {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 1rem;
		width: 100%;
		margin: 0 auto;
		box-sizing: border-box;
	}

	.top-bar-title {
		font-family: 'Caveat', cursive;
		font-size: 2rem;
	}
</style>
