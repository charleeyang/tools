import { ref } from 'vue'

const visible = ref(false)
const title = ref('')
const content = ref('')
const onOk = ref<(() => void) | null>(null)

export function useModal() {
  function openModal(t: string, body: string, ok?: () => void) {
    title.value = t
    content.value = body
    onOk.value = ok || null
    visible.value = true
  }
  function closeModal() {
    visible.value = false
  }
  return { visible, title, content, onOk, openModal, closeModal }
}
