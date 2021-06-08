<script>
	import { Blurb, Hero } from '@sveltejs/site-kit';
</script>


<svelte:head>
	<title>Progress.io - progress at your pace, follow your path.</title>

	<meta name="twitter:title" content="Svelte">
	<meta name="twitter:description" content="Cybernetically enhanced web apps">
	<meta name="Description" content="Cybernetically enhanced web apps">
</svelte:head>

<h1 class="visually-hidden">Svelte</h1>
<Hero
	title="Progress.io"
	tagline="Follow the path to become a master"
	outline="progress-logo-outline.svg"
	logotype="progress-logo.svg"
/>

<Blurb>
	<!-- TODO change href-->
	<a href="tutorial/go" slot="one">
		<h2>GO</h2>
		<p>The new most popular language from Google.</p>
		<span class="learn-more"></span>
	</a>

	<!-- TODO change href-->
	<a href="tutorial/c" slot="two">
		<h2>C</h2>
		<p>The most popular, powerful and flexible language ever.</p>

		<span class="learn-more"></span>
	</a>

	<!-- TODO change href-->
	<a href="tutorial/python" slot="three">
		<h2>Python</h2>
		<p>The most used language in data science, perfect for toy project or bigger one.</p>

		<span class="learn-more"></span>
	</a>

	<div class="description" slot="what">
		<p>Made by student at ESGI, Paris.</p>

		<p>With love, passion, coffee and LOTS of stress.</p>
	</div>
</Blurb>

