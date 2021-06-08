<script>
    import {setContext} from 'svelte';
    import {stores} from '@sapper/app';
    import {Icons, Nav, NavItem} from '@sveltejs/site-kit';

    export let segment;

    const {page, preloading, session} = stores();

    setContext('app', {
        login: () => {
            const login_window = window.open(`${window.location.origin}/auth/login`, 'login', 'width=600,height=400');

            window.addEventListener('message', function handler(event) {
                login_window.close();
                window.removeEventListener('message', handler);
                $session.user = event.data.user;
            });
        },

        logout: async () => {
            const r = await fetch(`/auth/logout`, {
                credentials: 'include'
            });

            if (r.ok) $session.user = null;
        }
    });
</script>

<Icons/>

<Nav {segment} {page} logo="progress-logo-horizontal.svg">
    <NavItem segment="tutorial">Login</NavItem>
    <NavItem segment="docs">Signup</NavItem>
</Nav>

<main>
    <slot></slot>
</main>

<style>
    main {
        position: relative;
        margin: 0 auto;
        /* padding: var(--nav-h) var(--side-nav) 0 var(--side-nav); */
        padding: var(--nav-h) 0 0 0;
        overflow-x: hidden;
    }
</style>
