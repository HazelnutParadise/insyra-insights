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
  export let texts: Record<string, string> = {};

  // 創建事件分發器
  const dispatch = createEventDispatcher();

  // 獲取實際使用的選項值（使用 i18n 翻譯作為預設值）
  $: actualOptions = {
    title: options.title || texts["dialog_defaults.alert_title"] || "提示",
    message: options.message,
    buttonText:
      options.buttonText || texts["dialog_defaults.confirm_button"] || "確定",
    type: options.type || "info",
  };

  // 處理確定按鈕點擊
  function handleOk() {
    dispatch("close", { action: "ok" });
  }
  // 處理背景點擊（已禁用，不允許點擊背景關閉）
  function handleBackdropClick() {
    // 不執行任何操作，禁止點擊背景關閉對話框
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
      class="alert-dialog {getThemeClass(actualOptions.type || 'info')}"
      on:click|stopPropagation
    >
      <div class="alert-header">
        <span class="alert-icon">{getIcon(actualOptions.type || "info")}</span>
        <h3 class="alert-title">{actualOptions.title}</h3>
      </div>

      <div class="alert-content">
        <p class="alert-message">{actualOptions.message}</p>
      </div>

      <div class="alert-footer">
        <button class="alert-button primary" on:click={handleOk} autofocus>
          {actualOptions.buttonText}
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
    background: rgba(0, 0, 0, 0.32); /* Material Design backdrop color */
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    animation: fadeIn var(--transition-duration) ease-out;
  }

  .alert-dialog {
    background: var(--surface-color); /* Material Design surface color */
    border-radius: var(
      --radius-small
    ); /* Material Design dialogs often have smaller border radius */
    box-shadow: var(--shadow-dialog); /* Use Material Design shadow */
    max-width: 420px;
    width: calc(
      100% - 2 * var(--spacing-lg)
    ); /* Ensure some margin on smaller screens */
    max-height: 80vh;
    overflow: hidden; /* Content will be scrollable if needed */
    animation: slideIn var(--transition-duration) cubic-bezier(0, 0, 0.2, 1); /* Material animation curve */
    display: flex;
    flex-direction: column;
  }

  .alert-header {
    display: flex;
    align-items: center;
    padding: var(--spacing-lg);
    /* Removed background gradient, Material Design dialogs usually have plain surface color */
  }

  .alert-icon {
    font-size: 24px; /* Adjusted icon size */
    margin-right: var(--spacing-md);
    color: var(
      --primary-color
    ); /* Icon color often matches primary color or is neutral */
    /* Removed background, border, and animation for a simpler Material look */
  }

  .alert-title {
    margin: 0;
    font-size: 1.25rem; /* Material Design title size (20px) */
    font-weight: 500; /* Material Design title weight */
    color: var(--text-primary);
    font-family: var(--font-primary); /* Ensure consistent font */
  }

  .alert-content {
    padding: 0 var(--spacing-lg) var(--spacing-lg); /* Adjust padding */
    flex-grow: 1; /* Allow content to take available space */
    overflow-y: auto; /* Make content scrollable if it overflows */
    color: var(--text-secondary);
    font-size: 1rem; /* Material Design body text (16px) */
    line-height: 1.5;
  }

  .alert-message {
    margin: 0;
  }

  .alert-footer {
    padding: var(--spacing-sm) var(--spacing-md); /* Reduced padding for Material spec */
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-sm);
    border-top: 1px solid rgba(0, 0, 0, 0.12); /* Material Design divider */
    /* Removed background gradient */
  }

  .alert-button {
    padding: var(--spacing-xs) var(--spacing-sm); /* Material text button padding */
    border: none;
    border-radius: var(--radius-small);
    font-size: 0.875rem; /* Material Design button text (14px) */
    font-weight: 500;
    cursor: pointer;
    transition: background-color var(--transition-fast);
    min-width: 64px; /* Material Design button min-width */
    text-transform: var(
      --text-button-text-transform
    ); /* Uppercase for Material buttons */
    letter-spacing: var(--text-button-letter-spacing);
    background-color: transparent; /* Text buttons are transparent */
    color: var(--primary-color); /* Text button color */
  }

  .alert-button:hover {
    background-color: rgba(
      var(--primary-color-rgb),
      0.08
    ); /* Slight background on hover */
  }

  .alert-button:active {
    background-color: rgba(
      var(--primary-color-rgb),
      0.12
    ); /* Slightly darker on active */
  }

  .alert-button:focus {
    outline: none;
    background-color: rgba(var(--primary-color-rgb), 0.12);
  }

  /* Theme variations - simplified for Material Design */
  /* Icons and specific colors can be adjusted further based on exact MD guidelines for success/warning/error */
  .alert-dialog.success .alert-icon {
    color: var(--success-color);
  }
  .alert-dialog.success .alert-button {
    color: var(--success-color);
  }
  .alert-dialog.success .alert-button:hover {
    background-color: rgba(76, 175, 80, 0.08); /* Success color RGB for hover */
  }
  .alert-dialog.success .alert-button:active,
  .alert-dialog.success .alert-button:focus {
    background-color: rgba(
      76,
      175,
      80,
      0.12
    ); /* Success color RGB for active/focus */
  }

  .alert-dialog.warning .alert-icon {
    color: var(--warning-color);
  }
  .alert-dialog.warning .alert-button {
    color: var(--warning-color);
  }
  .alert-dialog.warning .alert-button:hover {
    background-color: rgba(255, 152, 0, 0.08); /* Warning color RGB for hover */
  }
  .alert-dialog.warning .alert-button:active,
  .alert-dialog.warning .alert-button:focus {
    background-color: rgba(
      255,
      152,
      0,
      0.12
    ); /* Warning color RGB for active/focus */
  }

  .alert-dialog.error .alert-icon {
    color: var(--error-color);
  }
  .alert-dialog.error .alert-button {
    color: var(--error-color);
  }
  .alert-dialog.error .alert-button:hover {
    background-color: rgba(244, 67, 54, 0.08); /* Error color RGB for hover */
  }
  .alert-dialog.error .alert-button:active,
  .alert-dialog.error .alert-button:focus {
    background-color: rgba(
      244,
      67,
      54,
      0.12
    ); /* Error color RGB for active/focus */
  }

  /* Animations remain similar, but ensure Material easing curve */
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
      opacity: 0;
      transform: translateY(30px); /* Material dialogs often slide up */
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  /* Removed iconFloat animation as it's not typical for Material Design dialogs */

  /* Responsive design adjustments */
  @media (max-width: 480px) {
    .alert-dialog {
      margin: var(--spacing-medium);
      width: calc(100% - 2 * var(--spacing-medium));
    }

    .alert-header,
    .alert-content,
    .alert-footer {
      padding-left: var(--spacing-medium);
      padding-right: var(--spacing-medium);
    }
  }
</style>
