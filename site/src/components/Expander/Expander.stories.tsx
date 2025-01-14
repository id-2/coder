import { Expander } from "./Expander";
import type { Meta, StoryObj } from "@storybook/react";

const meta: Meta<typeof Expander> = {
  title: "components/Expander",
  component: Expander,
};

export default meta;
type Story = StoryObj<typeof Expander>;

export const Expanded: Story = {
  args: {
    expanded: true,
  },
};

export const Collapsed: Story = {
  args: {
    expanded: false,
  },
};
