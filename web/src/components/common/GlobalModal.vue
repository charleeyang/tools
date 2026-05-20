<template>
  <Teleport to="body">
    <div v-if="visible" class="modal-mask" @click.self="closeModal">
      <div class="modal-box">
        <div class="modal-head">
          <span>{{ title }}</span>
          <button class="close-btn" @click="closeModal">✕</button>
        </div>
        <div class="modal-body" v-html="content"></div>
        <div class="modal-foot">
          <button class="btn-cancel" @click="closeModal">取消</button>
          <button v-if="onOk" class="btn-primary" @click="onOk(); closeModal()">确定</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useModal } from '@/composables/useModal'
const { visible, title, content, onOk, closeModal } = useModal()
</script>

<style scoped>
.modal-mask {
  position: fixed; inset: 0; background: rgba(0,0,0,0.45); z-index: 9998;
  display: flex; align-items: center; justify-content: center;
}
.modal-box {
  background: #fff; border-radius: 8px; min-width: 420px; max-width: 600px;
  box-shadow: 0 6px 30px rgba(0,0,0,0.2);
}
.modal-head {
  display: flex; justify-content: space-between; align-items: center;
  padding: 16px 24px; border-bottom: 1px solid #f0f0f0; font-size: 16px; font-weight: 600;
}
.modal-body { padding: 24px; }
.modal-foot {
  display: flex; justify-content: flex-end; gap: 8px;
  padding: 12px 24px; border-top: 1px solid #f0f0f0;
}
.close-btn { background: none; border: none; font-size: 16px; cursor: pointer; color: #999; }
.btn-cancel { padding: 6px 16px; border: 1px solid #d9d9d9; border-radius: 6px; background: #fff; cursor: pointer; }
.btn-primary { padding: 6px 16px; border: none; border-radius: 6px; background: #1677ff; color: #fff; cursor: pointer; }
</style>
