<template>
  <!-- TODO: Сделать сами карточки -->
  <div class="grid grid-cols-4 grid-rows-3 gap-6">
    <NuxtLink
      to="/"
      class="blockBody"
      style="
        background-image: url('https://extxe.com/wp-content/uploads/2019/04/%D0%BB%D0%B5%D1%81-1.jpg');
      "
    />

    <NuxtLink class="blockBody col-span-2" />
    <NuxtLink class="blockBody row-span-2" />
    <NuxtLink class="blockBody row-span-2" />
    <NuxtLink class="blockBody col-span-2" />
    <NuxtLink class="blockBody" />
    <NuxtLink class="blockBody col-span-2" />
  </div>
</template>
