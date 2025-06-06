<script setup>
import VueDatePicker from "@vuepic/vue-datepicker"
import "@vuepic/vue-datepicker/dist/main.css"
import { useDateUtils } from "@/composables/useDateUtils"

const dateUtils = useDateUtils()
const date = defineModel("date", { required: true })
const props = defineProps({
  placeholder: { type: String, default: "" },
  required: { type: Boolean, default: false },
})
const formatter = (date) => dateUtils.dateToLocalDate(date)
</script>

<!-- TODO: enable month-picker option and set the value to the
last day of the selected month and set the time to 23:59:59 -->
<template>
  <VueDatePicker
    v-model="date"
    :placeholder="placeholder"
    :enableTimePicker="false"
    :required="required"
    :format="formatter"
    :max-date="new Date()"
    model-type="timestamp"
    input-class-name="input input-bordered w-full"
    popover-class-name="bg-base-100 border shadow-md rounded-lg z-[1000]"
  />
</template>