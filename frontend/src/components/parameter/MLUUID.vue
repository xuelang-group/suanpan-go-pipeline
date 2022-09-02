<template>
   <div class="param-panel-uuid">
     {{ this.uuid }}
   </div>
</template>

<script>
export default {
  name: 'MLUUID',
  props: {
    uuid: {
      required: true
    }
  },
}
</script>

<style>

</style>