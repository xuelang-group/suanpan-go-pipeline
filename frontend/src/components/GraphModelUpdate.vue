<template>
  <a-modal 
    v-model:visible="modalVisible"
    title="模型修改"
    :maskClosable="false"
    :footer="null"
    :afterClose="afterCloseHandler">
      <a-form
        ref="model"
        :model="modelState"
        :rules="rules"
        :label-col="labelCol"
        :wrapper-col="wrapperCol"
      >
      <a-form-item label="名称" name="name">
        <a-input v-model:value="modelState.name" />
      </a-form-item>
      <a-form-item label="描述" name="description">
        <a-input v-model:value="modelState.description" />
      </a-form-item>
      <a-form-item :wrapper-col="{ span: 4, offset: 19 }">
        <a-button type="primary" @click="onSubmit">修改</a-button>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script>
import { deepCopy } from '../utils/index'

export default {
  name: 'GraphModelUpdate',
  props: {
    model: {
      type: Object,
      required: true
    },
    visible: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:visible', 'confirm'],
  data() {
    return {
      modalVisible: this.visible,
      modelState: {
        name: this.model.name,
        description: this.model.description
      },
      rules: {
        name: [
           { required: true, message: '请输入模型名称', trigger: 'blur' }
        ]
      },
      labelCol: { span: 4 },
      wrapperCol: { span: 18 },
    }
  },
  watch: {
    visible() {
      this.modalVisible = this.visible
    }
  },
  methods: {
    afterCloseHandler() {
      this.$emit('update:visible', false)
    },
    onSubmit() {
      this.$refs.model.validate().then(() => {
        this.modalVisible = false
        let model = deepCopy(this.model)
        Object.assign(model, this.modelState)
        this.$emit('confirm', model)
      }).catch(() => {})
    },
  }
}
</script>

<style>

</style>