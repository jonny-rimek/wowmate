<template>
  <div id="app" class="antialiased text-gray-900 flex flex-col min-h-screen bg-gray-200">
	<nav x-data="{ open: true }" class="bg-white shadow">
	<div class="mx-auto px-2 sm:px-4 lg:px-8">
		<div class="flex justify-between h-16">
		<div class="flex px-2 lg:px-0">
			<div class="flex-shrink-0 flex items-center">
			<img class="block lg:hidden h-8 w-auto" src="/img/logos/workflow-mark-on-white.svg" alt="" />
			<img class="hidden lg:block h-8 w-auto" src="/img/logos/workflow-logo-on-white.svg" alt="" />
			</div>
			<div class="hidden lg:ml-6 lg:flex">
			<a href="#" class="inline-flex items-center px-1 pt-1 border-b-2 border-indigo-500 text-sm font-medium leading-5 text-gray-900 focus:outline-none focus:border-indigo-700 transition duration-150 ease-in-out">
				Ny'alotha 
			</a>
			<a href="#" class="ml-8 inline-flex items-center px-1 pt-1 border-b-2 border-transparent text-sm font-medium leading-5 text-gray-500 hover:text-gray-700 hover:border-gray-300 focus:outline-none focus:text-gray-700 focus:border-gray-300 transition duration-150 ease-in-out">
				Mythic Plus
			</a>
			<a href="#" class="ml-8 inline-flex items-center px-1 pt-1 border-b-2 border-transparent text-sm font-medium leading-5 text-gray-500 hover:text-gray-700 hover:border-gray-300 focus:outline-none focus:text-gray-700 focus:border-gray-300 transition duration-150 ease-in-out">
				Upload
			</a>
			</div>
		</div>
		<div class="flex-1 flex items-center justify-center px-2 lg:ml-6 lg:justify-end">
			<div class="max-w-lg w-full lg:max-w-xs">
			<label for="search" class="sr-only">Search</label>
			<div class="relative">
				<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
				<svg class="h-5 w-5 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
					<path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
				</svg>
				</div>
				<input id="search" class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:border-blue-300 focus:shadow-outline-blue sm:text-sm transition duration-150 ease-in-out" placeholder="Search" />
			</div>
			</div>
		</div>
		<div class="flex items-center lg:hidden">
			<button @click="open = !open" class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 focus:text-gray-500 transition duration-150 ease-in-out">
			<svg class="h-6 w-6" stroke="currentColor" fill="none" viewBox="0 0 24 24">
				<path :class="{'hidden': open, 'inline-flex': !open }" class="inline-flex" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
				<path :class="{'hidden': !open, 'inline-flex': open }" class="hidden" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
			</svg>
			</button>
		</div>
		<div class="hidden lg:ml-4 lg:flex lg:items-center">
			<button class="flex-shrink-0 p-1 border-2 border-transparent text-gray-400 rounded-full hover:text-gray-500 focus:outline-none focus:text-gray-500 focus:bg-gray-100 transition duration-150 ease-in-out">
			<svg class="h-6 w-6" stroke="currentColor" fill="none" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
			</svg>
			</button>
			<div @click.away="open = false" class="ml-4 relative flex-shrink-0" x-data="{ open: false }">
			<div>
				<button @click="open = !open" class="flex text-sm border-2 border-transparent rounded-full focus:outline-none focus:border-gray-300 transition duration-150 ease-in-out">
				<img class="h-8 w-8 rounded-full" src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80" alt="" />
				</button>
			</div>
			<div x-show="open" x-transition:enter="transition ease-out duration-100" x-transition:enter-start="transform opacity-0 scale-95" x-transition:enter-end="transform opacity-100 scale-100" x-transition:leave="transition ease-in duration-75" x-transition:leave-start="transform opacity-100 scale-100" x-transition:leave-end="transform opacity-0 scale-95" class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg">
				<div class="py-1 rounded-md bg-white shadow-xs">
				<a href="#" class="block px-4 py-2 text-sm leading-5 text-gray-700 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 transition duration-150 ease-in-out">Your Profile</a>
				<a href="#" class="block px-4 py-2 text-sm leading-5 text-gray-700 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 transition duration-150 ease-in-out">Settings</a>
				<a href="#" class="block px-4 py-2 text-sm leading-5 text-gray-700 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 transition duration-150 ease-in-out">Sign out</a>
				</div>
			</div>
			</div>
		</div>
		</div>
	</div>
	<div :class="{'block': open, 'hidden': !open}" class="hidden lg:hidden">
		<div class="pt-2 pb-3">
		<a href="#" class="block pl-3 pr-4 py-2 border-l-4 border-indigo-500 text-base font-medium text-indigo-700 bg-indigo-50 focus:outline-none focus:text-indigo-800 focus:bg-indigo-100 focus:border-indigo-700 transition duration-150 ease-in-out">Dashboard</a>
		<a href="#" class="mt-1 block pl-3 pr-4 py-2 border-l-4 border-transparent text-base font-medium text-gray-600 hover:text-gray-800 hover:bg-gray-50 hover:border-gray-300 focus:outline-none focus:text-gray-800 focus:bg-gray-50 focus:border-gray-300 transition duration-150 ease-in-out">Team</a>
		<a href="#" class="mt-1 block pl-3 pr-4 py-2 border-l-4 border-transparent text-base font-medium text-gray-600 hover:text-gray-800 hover:bg-gray-50 hover:border-gray-300 focus:outline-none focus:text-gray-800 focus:bg-gray-50 focus:border-gray-300 transition duration-150 ease-in-out">Projects</a>
		<a href="#" class="mt-1 block pl-3 pr-4 py-2 border-l-4 border-transparent text-base font-medium text-gray-600 hover:text-gray-800 hover:bg-gray-50 hover:border-gray-300 focus:outline-none focus:text-gray-800 focus:bg-gray-50 focus:border-gray-300 transition duration-150 ease-in-out">Calendar</a>
		</div>
		<div class="pt-4 pb-3 border-t border-gray-200">
		<div class="flex items-center px-4">
			<div class="flex-shrink-0">
			<img class="h-10 w-10 rounded-full" src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80" alt="" />
			</div>
			<div class="ml-3">
			<div class="text-base font-medium leading-6 text-gray-800">Tom Cook</div>
			<div class="text-sm font-medium leading-5 text-gray-500">tom@example.com</div>
			</div>
		</div>
		<div class="mt-3">
			<a href="#" class="block px-4 py-2 text-base font-medium text-gray-500 hover:text-gray-800 hover:bg-gray-100 focus:outline-none focus:text-gray-800 focus:bg-gray-100 transition duration-150 ease-in-out">Your Profile</a>
			<a href="#" class="mt-1 block px-4 py-2 text-base font-medium text-gray-500 hover:text-gray-800 hover:bg-gray-100 focus:outline-none focus:text-gray-800 focus:bg-gray-100 transition duration-150 ease-in-out">Settings</a>
			<a href="#" class="mt-1 block px-4 py-2 text-base font-medium text-gray-500 hover:text-gray-800 hover:bg-gray-100 focus:outline-none focus:text-gray-800 focus:bg-gray-100 transition duration-150 ease-in-out">Sign out</a>
		</div>
		</div>
	</div>
	</nav>

  <!-- Static sidebar for desktop -->
  <div class="hidden md:flex md:flex-shrink-0 h-screen">
    <div class="flex flex-col w-64 border-r border-gray-200 bg-white">
      <!-- <div class="flex items-center flex-shrink-0 px-4">
        <img class="h-8 w-auto" src="/img/logos/workflow-logo-on-white.svg" alt="Workflow" />
      </div> -->
      <div class="mt-2 h-0 flex-1 flex flex-col overflow-y-auto">
        <!-- Sidebar component, swap this element with another sidebar if you like -->
<nav class="flex-1 px-2 bg-white">
 
  <a href="#" class="group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-900 rounded-md bg-gray-100 hover:text-gray-900 hover:bg-gray-100 focus:outline-none focus:bg-gray-200 transition ease-in-out duration-150">
    <span class="truncate">
      Wrathion
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
		Maut
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
      Prophet Skitra
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
      Dark Inquisitor Xanesh
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
      The Hivemind
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
      Shad'har the Insatiable
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
      Drest'agath
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
      Vexiona
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
      Ra-den the Despoiled
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
      Il'gynoth, Corruption Reborn
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
      Carapace of N'Zoth
    </span>
  </a>
  <a href="#" class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-600 rounded-md hover:text-gray-900 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150">
    <span class="truncate">
      N'Zoth the Corruptor
    </span>
  </a>
</nav>
         </div>
    </div>
    <router-view class="flex-1"/>
    <div class="flex-shrink-0 py-12 text-center">
      <div>
        <router-link class="font-semibold text-gray-700 hover:underline" to="/dropdown">Dropdown</router-link>
        <span class="mx-4 text-gray-500">|</span>
        <router-link class="font-semibold text-gray-700 hover:underline" to="/modal">Modal</router-link>
      </div>
      <div class="mt-2">
        <a href="https://github.com/adamwathan/vue-tailwind-examples" class="text-sm font-semibold text-gray-600 hover:underline">
          Source on GitHub
        </a>
      </div>
    </div>
    <portal-target name="modals"></portal-target>
  </div>
  </div>
</template>

<style>
@tailwind base;
@tailwind components;
@tailwind utilities;

.origin-top-right {
  transform-origin: top right;
}
.transition-all {
  transition-property: all;
}
.transition-fastest {
  transition-duration: 50ms;
}
.transition-faster {
  transition-duration: 100ms;
}
.transition-fast {
  transition-duration: 150ms;
}
.transition-medium {
  transition-duration: 200ms;
}
.ease-out-quad {
  transition-timing-function: cubic-bezier(.25, .46, .45, .94);
}
.ease-in-quad {
  transition-timing-function: cubic-bezier(.55, .085, .68, .53);
}
.scale-70 {
  transform: scale(.7);
}
.scale-100 {
  transform: scale(1);
}
</style>
