<template>
  <button class="ml-button" :class="{ 'ml-button-disabled': disabled, 'ml-button-readonly': readonly }" :disabled="disabled">
    <i :class="['iconfont', iconfont]"></i>
    <span class="ml-button-label">{{ label }}</span>
  </button>
</template>

<script>
export default {
  name: 'MLButton',
  props: {
    iconfont: {
      type: String
    },
    label: {
      type: String,
      default: ''
    },
    disabled: {
      type: Boolean,
      default: false
    },
    readonly: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {}
  }
}
</script>

<style>

</style>