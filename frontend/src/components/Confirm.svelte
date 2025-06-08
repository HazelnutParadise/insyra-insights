<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { ConfirmOptions } from "../types/dialog";

  // 組件屬性
  export let visible: boolean = false;
  export let options: ConfirmOptions = {
    title: "確認",
    message: "",
    confirmText: "確定",
    cancelText: "取消",
    type: "info",
  };

  // 創建事件分發器
  const dispatch = createEventDispatcher();

  // 處理確認按鈕點擊
  function handleConfirm() {
    dispatch("close", { action: "confirm", result: true });
  }

  // 處理取消按鈕點擊
  function handleCancel() {
    dispatch("close", { action: "cancel", result: false });
  }
  // 處理背景點擊（已禁用，不允許點擊背景關閉）
  function handleBackdropClick() {
    // 不執行任何操作，禁止點擊背景關閉對話框
  }

  // 處理 ESC 鍵（等同於取消）
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      handleCancel();
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
    class="confirm-backdrop"
    on:click={handleBackdropClick}
    on:keydown={handleKeydown}
  >
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
      class="confirm-dialog {getThemeClass(options.type || 'info')}"
      on:click|stopPropagation
    >
      <div class="confirm-header">
        <span class="confirm-icon">{getIcon(options.type || "info")}</span>
        <h3 class="confirm-title">{options.title || "確認"}</h3>
      </div>

      <div class="confirm-content">
        <p class="confirm-message">{options.message}</p>
      </div>

      <div class="confirm-footer">
        <button class="confirm-button secondary" on:click={handleCancel}>
          {options.cancelText || "取消"}
        </button>
        <button
          class="confirm-button primary"
          on:click={handleConfirm}
          autofocus
        >
          {options.confirmText || "確定"}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .confirm-backdrop {
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

  .confirm-dialog {
    background: white;
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    max-width: 400px;
    width: 90%;
    max-height: 80vh;
    overflow: hidden;
    animation: slideIn 0.2s ease-out;
  }

  .confirm-header {
    display: flex;
    align-items: center;
    padding: 20px 20px 16px 20px;
    border-bottom: 1px solid #e0e0e0;
  }

  .confirm-icon {
    font-size: 24px;
    margin-right: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .confirm-title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #333;
  }

  .confirm-content {
    padding: 16px 20px;
  }

  .confirm-message {
    margin: 0;
    font-size: 14px;
    line-height: 1.5;
    color: #666;
  }

  .confirm-footer {
    padding: 16px 20px 20px 20px;
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    border-top: 1px solid #e0e0e0;
  }

  .confirm-button {
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
    min-width: 80px;
  }

  .confirm-button.secondary {
    background-color: #f5f5f5;
    color: #666;
    border: 1px solid #ddd;
  }

  .confirm-button.secondary:hover {
    background-color: #e0e0e0;
  }

  .confirm-button.primary {
    background-color: #2196f3;
    color: white;
  }

  .confirm-button.primary:hover {
    background-color: #1976d2;
  }
  .confirm-button.primary:focus,
  .confirm-button.secondary:focus {
    outline: none;
  }

  /* 主題變化 */
  .confirm-dialog.danger .confirm-header {
    border-bottom-color: #f44336;
  }

  .confirm-dialog.danger .confirm-button.primary {
    background-color: #f44336;
  }

  .confirm-dialog.danger .confirm-button.primary:hover {
    background-color: #d32f2f;
  }

  .confirm-dialog.warning .confirm-header {
    border-bottom-color: #ff9800;
  }

  .confirm-dialog.warning .confirm-button.primary {
    background-color: #ff9800;
  }

  .confirm-dialog.warning .confirm-button.primary:hover {
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
