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
    // type: "info", // type is not used in Material Design for confirm dialogs
  };
  export let texts: Record<string, string> = {};

  // 創建事件分發器
  const dispatch = createEventDispatcher();

  // 獲取實際使用的選項值（使用 i18n 翻譯作為預設值）
  $: actualOptions = {
    title: options.title || texts["dialog_defaults.confirm_title"] || "確認",
    message: options.message,
    confirmText:
      options.confirmText || texts["dialog_defaults.confirm_button"] || "確定",
    cancelText:
      options.cancelText || texts["dialog_defaults.cancel_button"] || "取消",
  };

  // 處理確認按鈕點擊
  function handleConfirm() {
    dispatch("close", { action: "confirm", result: true });
  }

  // 處理取消按鈕點擊
  function handleCancel() {
    dispatch("close", { action: "cancel", result: false });
  }

  // 處理背景點擊
  function handleBackdropClick() {
    // Material Design dialogs typically close on backdrop click
    // If this is not desired, this function can be left empty or call handleCancel()
    // For now, let's allow closing on backdrop click for consistency with typical MD behavior
    handleCancel();
  }

  // 處理 ESC 鍵（等同於取消）
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      handleCancel();
    }
  }
</script>

{#if visible}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div
    class="confirm-backdrop"
    on:click={handleBackdropClick}
    on:keydown={handleKeydown}
    role="dialog"
    aria-modal="true"
    aria-labelledby="confirm-title"
    aria-describedby="confirm-message"
  >
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="confirm-dialog" on:click|stopPropagation role="document">
      {#if actualOptions.title}
        <h2 class="confirm-title" id="confirm-title">{actualOptions.title}</h2>
      {/if}

      {#if actualOptions.message}
        <div class="confirm-content">
          <p class="confirm-message" id="confirm-message">
            {actualOptions.message}
          </p>
        </div>
      {/if}

      <div class="confirm-actions">
        <button class="confirm-button text-button" on:click={handleCancel}>
          {actualOptions.cancelText}
        </button>
        <button
          class="confirm-button text-button primary"
          on:click={handleConfirm}
          autofocus
        >
          {actualOptions.confirmText}
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
    background-color: rgba(0, 0, 0, 0.32); /* Material Design backdrop */
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    animation: dialog-fade-in var(--transition-standard) ease-out;
  }

  .confirm-dialog {
    background-color: var(--surface-color);
    border-radius: var(--radius-medium); /* Standard Material radius */
    box-shadow: var(--shadow-dialog); /* Material Design dialog shadow */
    min-width: 280px; /* Material Design minimum width */
    max-width: 560px; /* Material Design maximum width */
    width: calc(100% - 64px); /* Responsive width with padding */
    max-height: calc(100% - 64px);
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    animation: dialog-scale-open var(--transition-standard)
      cubic-bezier(0.4, 0, 0.2, 1);
    padding: 24px; /* Material Design padding */
    box-sizing: border-box;
  }

  .confirm-title {
    font-size: 1.25rem; /* 20sp */
    font-weight: 500;
    line-height: 1.6; /* 32sp */
    color: var(--text-primary);
    margin: 0 0 16px 0; /* Margin bottom if content exists */
    font-family: "Roboto", "Nunito", sans-serif; /* Material uses Roboto */
  }
  /* If there's no message, title has more bottom margin */
  .confirm-title:last-child {
    margin-bottom: 24px; /* Connects to actions */
  }

  .confirm-content {
    font-size: 1rem; /* 16sp */
    line-height: 1.5; /* 24sp */
    color: var(--text-secondary);
    margin: 0 0 24px 0; /* Space before actions */
    font-family: "Roboto", "Nunito", sans-serif;
  }
  .confirm-content p:first-child {
    margin-top: 0;
  }
  .confirm-content p:last-child {
    margin-bottom: 0;
  }
  .confirm-message {
    margin: 0; /* Reset paragraph margin */
  }

  .confirm-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px; /* Material Design gap for actions */
    padding-top: 8px; /* Add some space if content is short or no content */
  }

  .confirm-button {
    padding: 0 8px; /* Horizontal padding for text buttons */
    min-width: 64px; /* Minimum width for buttons */
    height: 36px; /* Standard height for text buttons */
    border: none;
    border-radius: var(--radius-small); /* 4dp */
    font-family: "Roboto", "Nunito", sans-serif;
    font-size: 0.875rem; /* 14sp */
    font-weight: 500;
    text-transform: var(--text-button-text-transform);
    letter-spacing: var(--text-button-letter-spacing);
    cursor: pointer;
    transition: background-color var(--transition-fast) ease-out;
    background-color: transparent;
    color: var(--primary-color); /* Default text button color */
    display: inline-flex;
    align-items: center;
    justify-content: center;
    box-sizing: border-box;
  }

  .confirm-button:hover {
    background-color: rgba(
      var(--primary-color-rgb),
      0.08
    ); /* Subtle hover for primary */
  }
  .confirm-button.text-button:not(.primary):hover {
    background-color: rgba(
      0,
      0,
      0,
      0.08
    ); /* Subtle hover for standard text button */
  }
  .confirm-button.text-button:not(.primary) {
    color: var(--text-secondary); /* Standard text button color */
  }
  .confirm-button.text-button:not(.primary):hover {
    background-color: rgba(
      var(--text-primary-rgb, 0, 0, 0),
      0.04
    ); /* Use text-primary-rgb if defined, else fallback */
  }

  .confirm-button:active {
    background-color: rgba(
      var(--primary-color-rgb),
      0.12
    ); /* Slightly darker for active */
  }
  .confirm-button.text-button:not(.primary):active {
    background-color: rgba(var(--text-primary-rgb, 0, 0, 0), 0.08);
  }

  .confirm-button.primary {
    color: var(--primary-color);
  }
  .confirm-button.primary:hover {
    background-color: rgba(var(--primary-color-rgb), 0.08);
  }
  .confirm-button.primary:active {
    background-color: rgba(var(--primary-color-rgb), 0.12);
  }

  .confirm-button:focus {
    outline: none;
    background-color: rgba(
      var(--primary-color-rgb),
      0.12
    ); /* Focus indicator */
  }
  .confirm-button.text-button:not(.primary):focus {
    background-color: rgba(var(--text-primary-rgb, 0, 0, 0), 0.08);
  }

  /* Animations */
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

  /* Responsive adjustments for smaller screens */
  @media (max-width: 600px) {
    .confirm-dialog {
      width: calc(100% - 32px); /* Smaller margins on mobile */
      max-width: calc(100% - 32px);
      border-radius: var(
        --radius-large
      ); /* Larger radius on mobile as per some MD guidelines */
      padding: 20px;
    }

    .confirm-title {
      font-size: 1.125rem; /* Slightly smaller title on mobile */
      margin-bottom: 16px;
    }
    .confirm-title:last-child {
      margin-bottom: 20px;
    }

    .confirm-content {
      font-size: 0.9375rem; /* Slightly smaller content text */
      margin-bottom: 20px;
    }

    .confirm-actions {
      padding-top: 0; /* No extra top padding if content is short */
    }
  }
</style>
