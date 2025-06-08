<script lang="ts">
  import { createEventDispatcher, onMount, tick } from "svelte";
  import type { InputOptions } from "../types/dialog";

  // 組件屬性
  export let visible: boolean = false;
  export let options: InputOptions = {
    title: "輸入",
    message: "",
    placeholder: "",
    defaultValue: "",
    confirmText: "確定",
    cancelText: "取消",
    inputType: "text",
    // type: "info", // type is not used in Material Design for input dialogs
  };
  export let texts: Record<string, string> = {};

  // 創建事件分發器
  const dispatch = createEventDispatcher();
  // 獲取實際使用的選項值（使用 i18n 翻譯作為預設值）
  $: actualOptions = {
    title: options.title || texts["dialog_defaults.input_title"] || "輸入",
    message: options.message,
    placeholder: options.placeholder,
    defaultValue: options.defaultValue,
    confirmText:
      options.confirmText || texts["dialog_defaults.confirm_button"] || "確定",
    cancelText:
      options.cancelText || texts["dialog_defaults.cancel_button"] || "取消",
    inputType: options.inputType || "text",
  };

  // 輸入值
  let inputValue = "";
  let inputElement: HTMLInputElement;
  let previousVisible = false; // 當對話框變為可見時，重置輸入值並設置焦點
  $: if (visible && !previousVisible) {
    // 總是重置為預設值，避免記憶性問題
    inputValue = actualOptions.defaultValue || "";
    tick().then(() => {
      if (inputElement) {
        inputElement.focus();
        if (inputValue) {
          inputElement.select();
        }
      }
    });
    previousVisible = visible;
  }

  // 當對話框變為不可見時，清空輸入值並更新 previousVisible
  $: if (!visible && previousVisible) {
    inputValue = "";
    previousVisible = visible;
  }

  $: previousVisible = visible;

  // 處理確認按鈕點擊
  function handleConfirm() {
    // console.log("Input 組件 - 確認按鈕點擊，inputValue:", inputValue);
    dispatch("close", { action: "confirm", result: inputValue });
  }

  // 處理取消按鈕點擊
  function handleCancel() {
    // console.log("Input 組件 - 取消按鈕點擊");
    dispatch("close", { action: "cancel", result: null });
  }

  // 處理背景點擊（關閉對話框）
  function handleBackdropClick() {
    handleCancel(); // Material Design 對話框通常在點擊背景時關閉
  }

  // 處理 ESC 鍵（等同於取消）和 Enter 鍵（在對話框上按 Enter 也應該確認）
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      handleCancel();
    }
    // 在對話框上（而不是在輸入框內）按 Enter 鍵也應該確認
    // 輸入元素有自己的 Enter 鍵處理程序
    if (event.key === "Enter" && event.target !== inputElement) {
      handleConfirm();
    }
  }

  // 處理輸入框的 Enter 鍵
  function handleInputKeydown(event: KeyboardEvent) {
    if (event.key === "Enter") {
      event.preventDefault(); // 如果在表單中，防止表單提交
      handleConfirm();
    }
  }
</script>

{#if visible}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div
    class="input-backdrop"
    on:click={handleBackdropClick}
    on:keydown={handleKeydown}
    role="dialog"
    aria-modal="true"
    aria-labelledby="input-title"
    aria-describedby={actualOptions.message ? "input-message" : undefined}
  >
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="input-dialog" on:click|stopPropagation role="document">
      {#if actualOptions.title}
        <h2 class="input-title" id="input-title">{actualOptions.title}</h2>
      {/if}
      <div class="input-content-area">
        {#if actualOptions.message}
          <p class="input-message" id="input-message">
            {actualOptions.message}
          </p>
        {/if}

        <!-- Material Design Text Field (Simplified) -->
        <div class="text-field">
          {#if actualOptions.inputType === "text"}
            <input
              type="text"
              class="input-field-md"
              placeholder=" "
              bind:value={inputValue}
              bind:this={inputElement}
              on:keydown={handleInputKeydown}
              aria-label={actualOptions.placeholder ||
                actualOptions.title ||
                "Input"}
              id="input-field-md-unique"
            />
          {:else if actualOptions.inputType === "password"}
            <input
              type="password"
              class="input-field-md"
              placeholder=" "
              bind:value={inputValue}
              bind:this={inputElement}
              on:keydown={handleInputKeydown}
              aria-label={actualOptions.placeholder ||
                actualOptions.title ||
                "Input"}
              id="input-field-md-unique"
            />
          {:else if actualOptions.inputType === "email"}
            <input
              type="email"
              class="input-field-md"
              placeholder=" "
              bind:value={inputValue}
              bind:this={inputElement}
              on:keydown={handleInputKeydown}
              aria-label={actualOptions.placeholder ||
                actualOptions.title ||
                "Input"}
              id="input-field-md-unique"
            />
          {:else if actualOptions.inputType === "number"}
            <input
              type="number"
              class="input-field-md"
              placeholder=" "
              bind:value={inputValue}
              bind:this={inputElement}
              on:keydown={handleInputKeydown}
              aria-label={actualOptions.placeholder ||
                actualOptions.title ||
                "Input"}
              id="input-field-md-unique"
            />
          {:else}
            <!-- Default to text if unknown -->
            <input
              type="text"
              class="input-field-md"
              placeholder=" "
              bind:value={inputValue}
              bind:this={inputElement}
              on:keydown={handleInputKeydown}
              aria-label={actualOptions.placeholder ||
                actualOptions.title ||
                "Input"}
              id="input-field-md-unique"
            />
          {/if}
          <label class="input-label-md" for="input-field-md-unique"
            >{actualOptions.placeholder ||
              actualOptions.title ||
              "Input"}</label
          >
        </div>
      </div>

      <div class="input-actions">
        <button class="input-button text-button" on:click={handleCancel}>
          {actualOptions.cancelText}
        </button>
        <button
          class="input-button text-button primary"
          on:click={handleConfirm}
        >
          {actualOptions.confirmText}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .input-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background-color: rgba(0, 0, 0, 0.32);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    animation: dialog-fade-in var(--transition-standard) ease-out;
  }

  .input-dialog {
    background-color: var(--surface-color);
    border-radius: var(--radius-medium);
    box-shadow: var(--shadow-dialog);
    min-width: 280px;
    max-width: 560px;
    width: calc(100% - 64px);
    max-height: calc(100% - 64px);
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    animation: dialog-scale-open var(--transition-standard)
      cubic-bezier(0.4, 0, 0.2, 1);
    padding: 24px;
    box-sizing: border-box;
  }

  .input-title {
    font-size: 1.25rem; /* 20sp */
    font-weight: 500;
    line-height: 1.6; /* 32sp */
    color: var(--text-primary);
    margin: 0 0 20px 0; /* Consistent margin */
    font-family: "Roboto", "Nunito", sans-serif;
  }

  .input-content-area {
    margin-bottom: 24px; /* Space before actions */
    flex-grow: 1;
  }

  .input-message {
    font-size: 1rem; /* 16sp */
    line-height: 1.5; /* 24sp */
    color: var(--text-secondary);
    margin: 0 0 16px 0; /* Space before input field */
    font-family: "Roboto", "Nunito", sans-serif;
  }
  .input-message:first-child {
    margin-top: 0;
  }
  .input-message:last-child {
    margin-bottom: 0; /* Remove margin if it's the last element before input */
  }

  /* Material Design Text Field Styles (Simplified) */
  .text-field {
    position: relative;
    padding-top: 16px; /* Space for the label to float */
    margin-bottom: 8px; /* Some space below the field */
  }

  .input-field-md {
    width: 100%;
    height: 56px; /* Standard height for MD text field */
    padding: 16px 12px 0; /* Padding for text, top padding for floating label */
    border: none;
    border-bottom: 1px solid rgba(0, 0, 0, 0.42); /* Resting line */
    border-radius: var(--radius-small) var(--radius-small) 0 0; /* Top corners rounded */
    font-family: "Roboto", "Nunito", sans-serif;
    font-size: 1rem; /* 16sp */
    background-color: var(
      --surface-variant
    ); /* Light background for the field */
    color: var(--text-primary);
    transition:
      border-bottom-color var(--transition-fast) ease-out,
      background-color var(--transition-fast) ease-out;
    box-sizing: border-box;
  }

  .input-field-md:hover {
    border-bottom-color: rgba(0, 0, 0, 0.87); /* Darker line on hover */
    background-color: #ececec; /* Slightly different background on hover */
  }

  .input-field-md:focus {
    outline: none;
    border-bottom: 2px solid var(--primary-color); /* Primary color line on focus */
    background-color: #e0e0e0; /* Slightly different background on focus */
  }

  .input-label-md {
    position: absolute;
    top: 32px; /* Vertically centered with input text before floating */
    left: 12px;
    font-size: 1rem; /* 16sp */
    color: var(--text-secondary);
    pointer-events: none;
    transition:
      transform var(--transition-fast) ease-out,
      font-size var(--transition-fast) ease-out,
      color var(--transition-fast) ease-out;
    transform-origin: left top;
  }

  /* Corrected floating label logic: 
     - Use :focus OR if the input has a value (NOT :placeholder-shown) 
     - The placeholder attribute on the input itself should be a single space for this to work reliably if no actual placeholder text is desired.
  */
  .input-field-md:focus + .input-label-md,
  .input-field-md:not(:placeholder-shown) + .input-label-md {
    transform: translateY(-16px) scale(0.75); /* Float label up and shrink */
    color: var(--primary-color);
  }

  /* Ensure placeholder is not shown when there is a value, to allow label to float correctly 
  .input-field-md:not(:placeholder-shown) {
     No specific style needed here, but :placeholder-shown is key 
  }
  */

  .input-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    padding-top: 8px;
  }

  .input-button {
    padding: 0 8px;
    min-width: 64px;
    height: 36px;
    border: none;
    border-radius: var(--radius-small);
    font-family: "Roboto", "Nunito", sans-serif;
    font-size: 0.875rem;
    font-weight: 500;
    text-transform: var(--text-button-text-transform);
    letter-spacing: var(--text-button-letter-spacing);
    cursor: pointer;
    transition: background-color var(--transition-fast) ease-out;
    background-color: transparent;
    color: var(--primary-color);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    box-sizing: border-box;
  }

  .input-button:hover {
    background-color: rgba(var(--primary-color-rgb), 0.08);
  }
  .input-button.text-button:not(.primary) {
    color: var(--text-secondary);
  }
  .input-button.text-button:not(.primary):hover {
    background-color: rgba(var(--text-primary-rgb, 0, 0, 0), 0.04);
  }

  .input-button:active {
    background-color: rgba(var(--primary-color-rgb), 0.12);
  }
  .input-button.text-button:not(.primary):active {
    background-color: rgba(var(--text-primary-rgb, 0, 0, 0), 0.08);
  }

  .input-button.primary {
    color: var(--primary-color);
  }

  .input-button:focus {
    outline: none;
    background-color: rgba(var(--primary-color-rgb), 0.12);
  }
  .input-button.text-button:not(.primary):focus {
    background-color: rgba(var(--text-primary-rgb, 0, 0, 0), 0.08);
  }

  /* 動畫 */
  @keyframes dialog-fade-in {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @keyframes dialog-scale-open {
    from {
      opacity: 0;
      transform: scale(0.95);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  /* 響應式調整 */
  @media (max-width: 600px) {
    .input-dialog {
      width: calc(100% - 32px);
      max-width: calc(100% - 32px);
      padding: 20px;
    }
    .input-title {
      font-size: 1.125rem;
      margin-bottom: 16px;
    }
    .input-content-area {
      margin-bottom: 20px;
    }
    .input-message {
      font-size: 0.9375rem;
      margin-bottom: 12px;
    }
    .text-field {
      padding-top: 12px;
    }
    .input-field-md {
      height: 52px;
      padding: 14px 10px 0;
    }
    .input-label-md {
      top: 28px;
      left: 10px;
      font-size: 0.9375rem;
    }
    .input-field-md:focus + .input-label-md,
    .input-field-md:not(:placeholder-shown) + .input-label-md {
      transform: translateY(-14px) scale(0.75);
    }
  }
</style>
