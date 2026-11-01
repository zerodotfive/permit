<script setup>
// POST target is the page the button is on (same path and query)
const action = location.pathname + location.search
</script>

<template>
  <AppLayout anon>
    <a-layout-content style="display: flex; align-items: center; justify-content: center">
      <!-- A real form submit: the browser POSTs to the current page and shows the response -->
      <form method="post" :action="action">
        <a-button type="primary" html-type="submit">Open</a-button>
      </form>
    </a-layout-content>
  </AppLayout>
</template>
