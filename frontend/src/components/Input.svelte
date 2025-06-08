<script lang="ts">
  import { createEventDispatcher } from "svelte";
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
    type: "info",
    inputType: "text",
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
    type: options.type || "info",
    inputType: options.inputType || "text",
  };

  // 輸入值
  let inputValue = "";
  // 當組件顯示時，重置輸入值為默認值
  $: if (visible) {
    inputValue = actualOptions.defaultValue || "";
  }

  // 處理確認按鈕點擊
  function handleConfirm() {
    dispatch("close", { action: "confirm", result: inputValue });
  }

  // 處理取消按鈕點擊
  function handleCancel() {
    dispatch("close", { action: "cancel", result: null });
  }
  // 處理背景點擊（已禁用，不允許點擊背景關閉）
  function handleBackdropClick() {
    // 不執行任何操作，禁止點擊背景關閉對話框
  }

  // 處理 ESC 鍵（等同於取消）
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      handleCancel();
    } else if (event.key === "Enter") {
      handleConfirm();
    }
  }

  // 處理輸入框的 Enter 鍵
  function handleInputKeydown(event: KeyboardEvent) {
    if (event.key === "Enter") {
      handleConfirm();
    }
  }

  // 獲取圖標
  function getIcon(type: string): string {
    switch (type) {
      case "danger":
        return "⚠️";
      case "warning":
        return "⚠️";
      case "info":
      default:
        return "❓";
    }
  }

  // 獲取主題類名
  function getThemeClass(type: string): string {
    switch (type) {
      case "danger":
        return "danger";
      case "warning":
        return "warning";
      case "info":
      default:
        return "info";
    }
  }
</script>

{#if visible}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div
    class="input-backdrop"
    on:click={handleBackdropClick}
    on:keydown={handleKeydown}
  >
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
      class="input-dialog {getThemeClass(actualOptions.type || 'info')}"
      on:click|stopPropagation
    >
      <div class="input-header">
        <span class="input-icon">{getIcon(actualOptions.type || "info")}</span>
        <h3 class="input-title">{actualOptions.title}</h3>
      </div>
      <div class="input-content">
        <p class="input-message">{actualOptions.message}</p>
        {#if actualOptions.inputType === "password"}
          <input
            type="password"
            class="input-field"
            placeholder={actualOptions.placeholder || ""}
            bind:value={inputValue}
            on:keydown={handleInputKeydown}
          />
        {:else if actualOptions.inputType === "email"}
          <input
            type="email"
            class="input-field"
            placeholder={actualOptions.placeholder || ""}
            bind:value={inputValue}
            on:keydown={handleInputKeydown}
          />
        {:else if actualOptions.inputType === "number"}
          <input
            type="number"
            class="input-field"
            placeholder={actualOptions.placeholder || ""}
            bind:value={inputValue}
            on:keydown={handleInputKeydown}
          />
        {:else}
          <input
            type="text"
            class="input-field"
            placeholder={actualOptions.placeholder || ""}
            bind:value={inputValue}
            on:keydown={handleInputKeydown}
          />
        {/if}
      </div>

      <div class="input-footer">
        <button class="input-button secondary" on:click={handleCancel}>
          {actualOptions.cancelText}
        </button>
        <button class="input-button primary" on:click={handleConfirm}>
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
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    animation: fadeIn 0.2s ease-out;
  }

  .input-dialog {
    background: white;
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    max-width: 450px;
    width: 90%;
    max-height: 80vh;
    overflow: hidden;
    animation: slideIn 0.2s ease-out;
  }

  .input-header {
    display: flex;
    align-items: center;
    padding: 20px 20px 16px 20px;
    border-bottom: 1px solid #e0e0e0;
  }

  .input-icon {
    font-size: 24px;
    margin-right: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .input-title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #333;
  }

  .input-content {
    padding: 16px 20px;
  }

  .input-message {
    margin: 0 0 16px 0;
    font-size: 14px;
    line-height: 1.5;
    color: #666;
  }

  .input-field {
    width: 100%;
    padding: 12px 16px;
    border: 2px solid #e0e0e0;
    border-radius: 4px;
    font-size: 14px;
    transition: border-color 0.2s;
    outline: none;
    box-sizing: border-box;
  }
  .input-field:focus {
    border-color: #2196f3;
    outline: none;
  }

  .input-footer {
    padding: 16px 20px 20px 20px;
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    border-top: 1px solid #e0e0e0;
  }

  .input-button {
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
    min-width: 80px;
  }

  .input-button.secondary {
    background-color: #f5f5f5;
    color: #666;
    border: 1px solid #ddd;
  }

  .input-button.secondary:hover {
    background-color: #e0e0e0;
  }

  .input-button.primary {
    background-color: #2196f3;
    color: white;
  }

  .input-button.primary:hover {
    background-color: #1976d2;
  }
  .input-button.primary:focus,
  .input-button.secondary:focus {
    outline: none;
  }

  /* 主題變化 */
  .input-dialog.danger .input-header {
    border-bottom-color: #f44336;
  }

  .input-dialog.danger .input-button.primary {
    background-color: #f44336;
  }

  .input-dialog.danger .input-button.primary:hover {
    background-color: #d32f2f;
  }

  .input-dialog.warning .input-header {
    border-bottom-color: #ff9800;
  }

  .input-dialog.warning .input-button.primary {
    background-color: #ff9800;
  }

  .input-dialog.warning .input-button.primary:hover {
    background-color: #f57c00;
  }

  /* 動畫 */
  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @keyframes slideIn {
    from {
      transform: translateY(-50px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }
</style>
