<template>
<div class="param-panel-form-vertical">
  <label class="param-panel-form-label">{{ paramInfo.name }}：</label>
   <div class="param-panel-group-inner">
    <div class="file-uploader-wrap">
      <a-input placeholder="请选择需要上传的文件" :title="currentValue" readonly :value="currentValue" :disabled="readonly" @click="fileUpload">
        <template v-slot:addonAfter>
          <a-button v-if="!paramInfo.loading" type="primary" @click="fileUpload" :disabled="readonly">文件上传</a-button>
          <a-button v-if="paramInfo.loading" type="primary" :disabled="readonly">
            <template #icon><LoadingOutlined /></template>上传中
          </a-button>
        </template>
      </a-input>
      <input ref="file" type="file" style="display:none" @change="fileHandler" />
    </div>
   </div>
</div>
<div v-show="showError" class="param-validate-tip">{{ errorMsg }}</div>
</template>

<script>
import { ossService } from '../../service'
import { requiredValidate } from '../../utils/validate'

export default {
  name: 'MLFileUploader',
  props: {
    paramInfo: {
      type: Object,
      required: true
    },
    params: {
      type: Array,
      required: true
    },
    readonly: {
      type: Boolean
    }
  },
  emits: ['param-change', 'valid-change'],
  data() {
    return {
      currentValue: this.paramInfo.value,
      showError: !this.paramInfo.valid,
      errorMsg: this.paramInfo.errMsg,
    }
  },
  watch: {
    showError(val) {
      this.paramInfo.valid = !val
      this.$emit('valid-change')
    },
    errorMsg(val) {
      this.paramInfo.errMsg = val
    },
    currentValue() {
      this.checkValid(this.currentValue)
      this.$emit('param-change', this.currentValue, this.paramInfo, this)
    }
  },
  created() {
    this.$parent.paramCompInsts.push(this); 
  },
  beforeUnmount() {
    if(this.$parent.paramCompInsts) {
      let idx = this.$parent.paramCompInsts.indexOf(this)
      if(idx > -1) {
        this.$parent.paramCompInsts.splice(idx, 1)
      }
    }
  },
  methods: {
    checkValid(val) {
      if(!requiredValidate(val, this.paramInfo, this.params)) {
        this.showError = true
        this.errorMsg = '该参数必填'
      }else {
        this.showError = false
        this.errorMsg = ''
      }
    },
    triggerValidCheck() {
      this.checkValid(this.currentValue)
    },
    fileUpload() {
      if(this.paramInfo.loading) {
        return
      }
      this.$refs.file.click()
    },
    fileHandler(e) {
      if(e.target && e.target.files && (e.target.files.length < 1)) {
        return;
      }
      let file = e.target.files[0]
      e.target.value = null
      this.paramInfo.loading = true
      let key = `${this.$store.state.storageNodePath}/${file.name}`
      ossService.upload(
        file, 
        key, 
        () => {},
        (err) => {
          this.paramInfo.loading = false
          console.error(`upload file error, key:${key}`, err)
        },
        () => {
          console.log(`upload file success, key:${key}`)
          this.paramInfo.loading = false
          this.currentValue = key
        }
        );
    }
  }
}
</script>

<style>

</style>