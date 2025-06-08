<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { AlertOptions } from "../types/dialog";

  // 組件屬性
  export let visible: boolean = false;
  export let options: AlertOptions = {
    title: "提示",
    message: "",
    buttonText: "確定",
    type: "info",
  };

  // 創建事件分發器
  const dispatch = createEventDispatcher();

  // 處理確定按鈕點擊
  function handleOk() {
    dispatch("close", { action: "ok" });
  }

  // 處理背景點擊（可選，是否允許點擊背景關閉）
  function handleBackdropClick() {
    // 可以選擇是否允許點擊背景關閉 Alert
    // handleOk();
  }

  // 處理 ESC 鍵
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      handleOk();
    }
  }

  // 獲取圖標
  function getIcon(type: string): string {
    switch (type) {
      case "success":
        return "✅";
      case "warning":
        return "⚠️";
      case "error":
        return "❌";
      case "info":
      default:
        return "ℹ️";
    }
  }

  // 獲取主題類名
  function getThemeClass(type: string): string {
    switch (type) {
      case "success":
        return "success";
      case "warning":
        return "warning";
      case "error":
        return "error";
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
    class="alert-backdrop"
    on:click={handleBackdropClick}
    on:keydown={handleKeydown}
  >
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
      class="alert-dialog {getThemeClass(options.type || 'info')}"
      on:click|stopPropagation
    >
      <div class="alert-header">
        <span class="alert-icon">{getIcon(options.type || "info")}</span>
        <h3 class="alert-title">{options.title || "提示"}</h3>
      </div>

      <div class="alert-content">
        <p class="alert-message">{options.message}</p>
      </div>

      <div class="alert-footer">
        <button class="alert-button primary" on:click={handleOk} autofocus>
          {options.buttonText || "確定"}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .alert-backdrop {
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

  .alert-dialog {
    background: white;
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    max-width: 400px;
    width: 90%;
    max-height: 80vh;
    overflow: hidden;
    animation: slideIn 0.2s ease-out;
  }

  .alert-header {
    display: flex;
    align-items: center;
    padding: 20px 20px 16px 20px;
    border-bottom: 1px solid #e0e0e0;
  }

  .alert-icon {
    font-size: 24px;
    margin-right: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .alert-title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #333;
  }

  .alert-content {
    padding: 16px 20px;
  }

  .alert-message {
    margin: 0;
    font-size: 14px;
    line-height: 1.5;
    color: #666;
  }

  .alert-footer {
    padding: 16px 20px 20px 20px;
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    border-top: 1px solid #e0e0e0;
  }

  .alert-button {
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
    min-width: 80px;
  }

  .alert-button.primary {
    background-color: #2196f3;
    color: white;
  }
  .alert-button.primary:hover {
    background-color: #1976d2;
  }

  .alert-button.primary:focus {
    outline: none;
  }

  /* 主題變化 */
  .alert-dialog.success .alert-header {
    border-bottom-color: #4caf50;
  }

  .alert-dialog.success .alert-button.primary {
    background-color: #4caf50;
  }

  .alert-dialog.success .alert-button.primary:hover {
    background-color: #388e3c;
  }

  .alert-dialog.warning .alert-header {
    border-bottom-color: #ff9800;
  }

  .alert-dialog.warning .alert-button.primary {
    background-color: #ff9800;
  }

  .alert-dialog.warning .alert-button.primary:hover {
    background-color: #f57c00;
  }

  .alert-dialog.error .alert-header {
    border-bottom-color: #f44336;
  }

  .alert-dialog.error .alert-button.primary {
    background-color: #f44336;
  }

  .alert-dialog.error .alert-button.primary:hover {
    background-color: #d32f2f;
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
